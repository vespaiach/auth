package mysqlrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
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
func (r *MysqlRoleActionRepo) GetByID(id int64) (*model.RoleAction, error) {
	rows, err := r.DbClient.Queryx(sqlGetRoleActionByID, id)
	if err != nil {
		log.Error("MysqlRoleActionRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	roleAction, err := mapRoleActionRow(rows)
	if err != nil {
		log.Error("MysqlRoleActionRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return roleAction, nil
}

var sqlCreateRoleAction = `
INSERT INTO role_actions(role_id, action_id) VALUES(?, ?);
`

// Create a new role-action
func (r *MysqlRoleActionRepo) Create(roleID int64, actionID int64) (int64, error) {
	stmt, err := r.DbClient.Prepare(sqlCreateRoleAction)
	if err != nil {
		log.Error("MysqlRoleActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(roleID, actionID)
	if err != nil {
		log.Error("MysqlRoleActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlRoleActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	return lastID, nil
}

var sqlDeleteRoleAction = `
DELETE FROM role_actions WHERE role_actions.id = ?;
`

// Delete role-action
func (r *MysqlRoleActionRepo) Delete(id int64) error {
	stmt, err := r.DbClient.Prepare(sqlDeleteRoleAction)
	if err != nil {
		log.Error("MysqlRoleActionRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	res, err := stmt.Exec(id)
	if err != nil {
		log.Error("MysqlRoleActionRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		log.Error("MysqlRoleActionRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
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
LIMIT :offset, :limit;`

const sqlCountListRoleAction = `
SELECT Count(*)
FROM roles INNER JOIN role_actions
ON roles.id = role_actions.role_id
INNER JOIN actions
ON role_actions.action_id = actions.id
%s ;`

// Query a list of role-actions
func (r *MysqlRoleActionRepo) Query(page int, perPage int, filters map[string]interface{}) ([]*model.RoleAction, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListRoleAction, conditions), filters)
		if err != nil {
			log.Error("MysqlRoleActionRepo - Query:", err)
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

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListRoleAction, conditions), filters)
	if err != nil {
		log.Error("MysqlRoleActionRepo - Query:", err)
		fmt.Println(fmt.Sprintf(sqlListRoleAction, conditions))
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.RoleAction, 0, perPage)
	for rows.Next() {
		ac, err := mapRoleActionRow(rows)
		if err != nil {
			log.Error("MysqlRoleActionRepo - Query:", err)
			return nil, 0, comtype.ErrQueryDataFailed
		}
		results = append(results, ac)
	}

	total := <-ch
	if total == -1 {
		return nil, 0, comtype.ErrQueryDataFailed
	}

	return results, total, nil
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
