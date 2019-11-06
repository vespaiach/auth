package mysqlrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlRoleActionRepo will implement model.RoleActionRepo
type MysqlRoleActionRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlRoleActionRepo create new instance of MysqlRoleActionRepo
func NewMysqlRoleActionRepo(db *sqlx.DB) model.RoleActionRepo {
	return &MysqlRoleActionRepo{
		db,
	}
}

var sqlGetRoleActionByID = `
SELECT
	role_actions.id,
	role_actions.role_id,
	role_actions.action_id,
	role_actions.created_at,
	actions.id,
	actions.action_name,
	actions.action_desc,
	actions.active,
	actions.created_at,
	actions.updated_at,
	roles.id,
	roles.role_name,
	roles.role_desc,
	roles.active,
	roles.created_at,
	roles.updated_at
FROM roles INNER JOIN role_actions
ON roles.id = role_actions.role_id
INNER JOIN actions
ON role_actions.action_id = actions.id
WHERE role_actions.id =?
LIMIT 1;
`

// GetByID find a role-action by its ID
func (r *MysqlRoleActionRepo) GetByID(id int64) (*model.RoleAction, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetRoleActionByID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleActionRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.NewCommonError(nil, "MysqlRoleActionRepo - GetByID:", comtype.ErrDataNotFound, nil)
	}

	roleAction, err := mapRoleActionRow(rows)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleActionRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}

	return roleAction, nil
}

var sqlCreateRoleAction = `
INSERT INTO role_actions(role_id, action_id) VALUES(?, ?);
`

// Create a new role-action
func (r *MysqlRoleActionRepo) Create(roleID int64, actionID int64) (int64, *comtype.CommonError) {
	stmt, err := r.DbClient.Prepare(sqlCreateRoleAction)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlRoleActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(roleID, actionID)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlRoleActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlRoleActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return lastID, nil
}

var sqlDeleteRoleAction = `
DELETE FROM role_actions WHERE role_actions.id = ?;
`

// Delete role-action
func (r *MysqlRoleActionRepo) Delete(id int64) *comtype.CommonError {
	stmt, err := r.DbClient.Prepare(sqlDeleteRoleAction)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlRoleActionRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlRoleActionRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		return comtype.NewCommonError(err, "MysqlRoleActionRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	return nil
}

const sqlListRoleAction = `
SELECT
	role_actions.id,
	role_actions.role_id,
	role_actions.action_id,
	role_actions.created_at,
	actions.id,
	actions.action_name,
	actions.action_desc,
	actions.active,
	actions.created_at,
	actions.updated_at,
	roles.id,
	roles.role_name,
	roles.role_desc,
	roles.active,
	roles.created_at,
	roles.updated_at
FROM roles INNER JOIN role_actions
ON roles.id = role_actions.role_id
INNER JOIN actions
ON role_actions.action_id = actions.id
%s
ORDER BY role_actions.created_at DESC
LIMIT :limit;`

// Query a list of role-actions
func (r *MysqlRoleActionRepo) Query(take int, filters map[string]interface{}) ([]*model.RoleAction, *comtype.CommonError) {
	conditions := sqlWhereBuilder(" AND ", filters)
	filters = sqlLikeConditionFilter(filters)
	if take == 0 {
		filters["limit"] = 100
	} else {
		filters["limit"] = take
	}

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListRoleAction, conditions), filters)
	if err != nil {
		fmt.Println(fmt.Sprintf(sqlListRoleAction, conditions))
		return nil, comtype.NewCommonError(err, "MysqlRoleActionRepo - Query:", comtype.ErrHandleDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.RoleAction, 0, take)
	for rows.Next() {
		ac, err := mapRoleActionRow(rows)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlRoleActionRepo - Query:", comtype.ErrHandleDataFail, nil)
		}
		results = append(results, ac)
	}

	if rows.Err() != nil {
		return nil, comtype.NewCommonError(rows.Err(), "MysqlRoleActionRepo - Query:", comtype.ErrHandleDataFail, nil)
	}

	return results, nil
}

func mapRoleActionRow(rows *sqlx.Rows) (*model.RoleAction, error) {
	var (
		ra model.RoleAction
		a  model.Action
		r  model.Role
	)

	err := rows.Scan(&ra.ID, &ra.RoleID, &ra.ActionID, &ra.CreatedAt,
		&a.ID, &a.ActionName, &a.ActionDesc, &a.Active, &a.CreatedAt, &a.UpdatedAt,
		&r.ID, &r.RoleName, &r.RoleDesc, &r.Active, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ra.Action = &a
	ra.Role = &r

	return &ra, nil
}
