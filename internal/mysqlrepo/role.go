package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlRoleRepo will implement model.RoleRepo
type MysqlRoleRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlRoleRepo create new instance of MysqlRoleRepo
func NewMysqlRoleRepo(db *sqlx.DB) model.RoleRepo {
	return &MysqlRoleRepo{
		db,
	}
}

var sqlGetRoleByID = `
SELECT *
FROM roles
WHERE roles.id =?
LIMIT 1;
`

// GetByID find a role by its ID
func (r *MysqlRoleRepo) GetByID(id int64) (*model.Role, error) {
	rows, err := r.DbClient.Queryx(sqlGetRoleByID, id)
	if err != nil {
		log.Error("MysqlRoleRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	role := new(model.Role)
	err = rows.StructScan(role)
	if err != nil {
		log.Error("MysqlRoleRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return role, nil
}

var sqlGetRoleByName = `
SELECT *
FROM roles
WHERE roles.role_name = ?
LIMIT 1;
`

// GetByName find a role by its name
func (r *MysqlRoleRepo) GetByName(name string) (*model.Role, error) {
	rows, err := r.DbClient.Queryx(sqlGetRoleByName, name)
	if err != nil {
		log.Error("MysqlRoleRepo - GetByName:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	role := new(model.Role)
	err = rows.StructScan(role)
	if err != nil {
		log.Error("MysqlRoleRepo - GetByName:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return role, nil
}

var sqlCreateRole = `
INSERT INTO roles(role_name, role_desc) VALUES(?, ?);
`

// Create a new role
func (r *MysqlRoleRepo) Create(name string, desc string) (int64, error) {
	stmt, err := r.DbClient.Prepare(sqlCreateRole)
	if err != nil {
		log.Error("MysqlRoleRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(name, desc)
	if err != nil {
		log.Error("MysqlRoleRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlRoleRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	return lastID, nil
}

var sqlUpdateRole = `
UPDATE roles
%s
WHERE roles.id = :id;
`

// Update role
func (r *MysqlRoleRepo) Update(id int64, fields map[string]interface{}) error {
	conditions := sqlWhereBuilder(", ", fields)
	if len(conditions) == 0 {
		log.Error("MysqlRoleRepo - Update:", errors.New("empty updating fields"))
		return comtype.ErrUpdateDataFailed
	}

	_, err := r.DbClient.NamedExec(fmt.Sprintf(sqlUpdateRole, conditions), fields)
	if err != nil {
		log.Error("MysqlRoleRepo - Update:", err)
		return comtype.ErrUpdateDataFailed
	}

	return nil
}

const sqlListRole = `
SELECT *
FROM roles
%s
ORDER BY %s 
LIMIT :offset, :limit;`

const sqlCountListRole = `
SELECT Count(*)
FROM roles
%s ;`

// Query a list of roles
func (r *MysqlRoleRepo) Query(page int64, perPage int64, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.Role, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListRole, conditions), filters)
		if err != nil {
			log.Error("MysqlRoleRepo - Query:", err)
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

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListRole, conditions, sortings), filters)
	if err != nil {
		log.Error("MysqlRoleRepo - Query:", err)
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.Role, 0, perPage)
	for rows.Next() {
		var ac model.Role
		rows.StructScan(&ac)
		if err != nil {
			log.Error("MysqlRoleRepo - Query:", err)
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
