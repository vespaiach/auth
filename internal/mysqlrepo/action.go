package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlActionRepo will implement model.ActionRepo
type MysqlActionRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlActionRepo create new instance of MysqlActionRepo
func NewMysqlActionRepo(db *sqlx.DB) model.ActionRepo {
	return &MysqlActionRepo{
		db,
	}
}

var sqlGetActionByID = `
SELECT 
	id,
	action_name,
	action_desc,
	active,
	created_at,
	updated_at
FROM actions
WHERE actions.id =?
LIMIT 1;
`

// GetByID find an action by its ID
func (r *MysqlActionRepo) GetByID(id int64) (*model.Action, error) {
	rows, err := r.DbClient.Queryx(sqlGetActionByID, id)
	if err != nil {
		log.Error("MysqlActionRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	action := new(model.Action)
	err = rows.StructScan(action)
	if err != nil {
		log.Error("MysqlActionRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return action, nil
}

var sqlGetActionByName = `
SELECT *
FROM actions
WHERE actions.action_name = ?
LIMIT 1;
`

// GetByName find an action by its name
func (r *MysqlActionRepo) GetByName(name string) (*model.Action, error) {
	rows, err := r.DbClient.Queryx(sqlGetActionByName, name)
	if err != nil {
		log.Error("MysqlActionRepo - GetByName:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	action := new(model.Action)
	err = rows.StructScan(action)
	if err != nil {
		log.Error("MysqlActionRepo - GetByName:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return action, nil
}

var sqlCreateAction = `
INSERT INTO actions(action_name, action_desc) VALUES(?, ?);
`

// Create an new action
func (r *MysqlActionRepo) Create(name string, desc string) (int64, error) {
	stmt, err := r.DbClient.Prepare(sqlCreateAction)
	if err != nil {
		log.Error("MysqlActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(name, desc)
	if err != nil {
		log.Error("MysqlActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	return lastID, nil
}

var sqlUpdateAction = `
UPDATE actions
%s
WHERE actions.id = :id;
`

// Update action
func (r *MysqlActionRepo) Update(id int64, fields map[string]interface{}) error {
	conditions := sqlWhereBuilder(", ", fields)
	if len(conditions) == 0 {
		log.Error("MysqlActionRepo - Update:", errors.New("empty updating fields"))
		return comtype.ErrUpdateDataFailed
	}

	_, err := r.DbClient.NamedExec(fmt.Sprintf(sqlUpdateAction, conditions), fields)
	if err != nil {
		log.Error("MysqlActionRepo - Update:", err)
		return comtype.ErrUpdateDataFailed
	}

	return nil
}

const sqlListAction = `
SELECT *
FROM actions
%s
ORDER BY %s 
LIMIT :offset, :limit;`

const sqlCountListAction = `
SELECT Count(*)
FROM actions
%s ;`

// Query a list of actions
func (r *MysqlActionRepo) Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.Action, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListAction, conditions), filters)
		if err != nil {
			log.Error("MysqlActionRepo - Query:", err)
			ch <- int64(-1)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			ch <- int64(-1)
			return
		}

		rows.Scan(&totals)
		ch <- totals
		close(ch)
	}()

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListAction, conditions, sortings), filters)
	if err != nil {
		log.Error("MysqlActionRepo - Query:", err)
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.Action, 0, perPage)
	for rows.Next() {
		var ac model.Action
		rows.StructScan(&ac)
		if err != nil {
			log.Error("MysqlActionRepo - Query:", err)
			return nil, 0, comtype.ErrQueryDataFailed
		}
		results = append(results, &ac)
	}

	total := <-ch
	if total == -1 {
		return nil, 0, comtype.ErrQueryDataFailed
	}

	return results, total, nil
}
