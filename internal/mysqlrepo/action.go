package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

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
func (r *MysqlActionRepo) GetByID(id int64) (*model.Action, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetActionByID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	action := new(model.Action)
	err = rows.StructScan(action)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByID:", comtype.ErrQueryDataFail, nil)
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
func (r *MysqlActionRepo) GetByName(name string) (*model.Action, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetActionByName, name)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByName:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	action := new(model.Action)
	err = rows.StructScan(action)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByName:", comtype.ErrQueryDataFail, nil)
	}

	return action, nil
}

var sqlCreateAction = `
INSERT INTO actions(action_name, action_desc) VALUES(?, ?);
`

// Create an new action
func (r *MysqlActionRepo) Create(name string, desc string) (int64, *comtype.CommonError) {
	stmt, err := r.DbClient.Prepare(sqlCreateAction)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(name, desc)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return lastID, nil
}

var sqlUpdateAction = `
UPDATE actions 
SET 
%s 
WHERE actions.id = :id;
`

// Update action
func (r *MysqlActionRepo) Update(id int64, fields map[string]interface{}) *comtype.CommonError {
	if len(fields) == 0 {
		return comtype.NewCommonError(errors.New("empty updating fields"), "MysqlActionRepo - Update:",
			comtype.ErrHandleDataFail, nil)
	}

	fields["id"] = id
	_, err := r.DbClient.NamedExec(fmt.Sprintf(sqlUpdateAction, sqlUpdateBuilder(fields, map[string]bool{"id": true})), fields)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlActionRepo - Update:", comtype.ErrHandleDataFail, nil)
	}

	return nil
}

const sqlListAction = `
SELECT *
FROM actions
%s
ORDER BY %s 
LIMIT :limit;`

// Query a list of actions
func (r *MysqlActionRepo) Query(take int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.Action, *comtype.CommonError) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	if take == 0 {
		filters["limit"] = 100
	} else {
		filters["limit"] = take
	}

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListAction, conditions, sortings), filters)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - Query:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.Action, 0, take)
	for rows.Next() {
		var ac model.Action
		rows.StructScan(&ac)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlActionRepo - Query:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, &ac)
	}

	if rows.Err() != nil {
		return nil, comtype.NewCommonError(rows.Err(), "MysqlActionRepo - Query:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
}

var sqlGetActionByUserID = `
(SELECT DISTINCT
	actions.id,
	actions.action_name,
	actions.action_desc,
	actions.active,
	actions.created_at,
	actions.updated_at
FROM user_roles INNER JOIN roles
ON user_roles.role_id = roles.id
INNER JOIN role_actions
ON user_roles.role_id = role_actions.role_id
INNER JOIN actions
ON role_actions.action_id = actions.id
INNER JOIN user_actions
ON user_actions.user_id = user_roles.user_id
WHERE
	user_roles.user_id = ? AND
	actions.active = 1 AND
	roles.active = 1)

UNION DISTINCT

(SELECT
	actions.id,
	actions.action_name,
	actions.action_desc,
	actions.active,
	actions.created_at,
	actions.updated_at
FROM users INNER JOIN user_actions
ON users.id = user_actions.id
INNER JOIN actions
ON actions.id = user_actions.action_id
WHERE
	users.id = ? AND
	actions.active = 1)
`

// GetByUserID gets list of action by user ID
func (r *MysqlActionRepo) GetByUserID(userID int64) ([]*model.Action, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetActionByUserID, userID, userID)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.Action, 0)
	for rows.Next() {
		var ac model.Action
		rows.StructScan(&ac)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, &ac)
	}

	err = rows.Err()
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlActionRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
}
