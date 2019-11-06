package mysqlrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlUserActionRepo will implement model.UserActionRepo
type MysqlUserActionRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlUserActionRepo create new instance of MysqlUserActionRepo
func NewMysqlUserActionRepo(db *sqlx.DB) model.UserActionRepo {
	return &MysqlUserActionRepo{
		db,
	}
}

var sqlGetUserActionByID = `
SELECT
	user_actions.id,
	user_actions.user_id,
	user_actions.action_id,
	user_actions.created_at,
	users.id,
	users.full_name,
	users.username,
	users.email,
	users.verified,
	users.active,
	users.created_at,
	users.updated_at,
	actions.id,
	actions.action_name,
	actions.action_desc,
	actions.active,
	actions.created_at,
	actions.updated_at
FROM users INNER JOIN user_actions
ON users.id = user_actions.user_id
INNER JOIN actions
ON user_actions.action_id = actions.id
WHERE user_actions.id =?
LIMIT 1;
`

// GetByID find a user-action by its ID
func (r *MysqlUserActionRepo) GetByID(id int64) (*model.UserAction, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetUserActionByID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserActionRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.NewCommonError(nil, "MysqlUserActionRepo - GetByID:", comtype.ErrDataNotFound, nil)
	}

	userAction, err := mapUserActionRow(rows)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserActionRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}

	return userAction, nil
}

var sqlCreateUserAction = `
INSERT INTO user_actions(user_id, action_id) VALUES(?, ?);
`

// Create a new user-action
func (r *MysqlUserActionRepo) Create(userID int64, actionID int64) (int64, *comtype.CommonError) {
	stmt, err := r.DbClient.Prepare(sqlCreateUserAction)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(userID, actionID)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserActionRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return lastID, nil
}

var sqlDeleteUserAction = `
DELETE FROM user_actions WHERE user_actions.id = ?;
`

// Delete user-action
func (r *MysqlUserActionRepo) Delete(id int64) *comtype.CommonError {
	stmt, err := r.DbClient.Prepare(sqlDeleteUserAction)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlUserActionRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlUserActionRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		return comtype.NewCommonError(err, "MysqlUserActionRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	return nil
}

const sqlListUserAction = `
SELECT
	user_actions.id,
	user_actions.user_id,
	user_actions.action_id,
	user_actions.created_at,
	users.id,
	users.full_name,
	users.username,
	users.email,
	users.verified,
	users.active,
	users.created_at,
	users.updated_at,
	actions.id,
	actions.action_name,
	actions.action_desc,
	actions.active,
	actions.created_at,
	actions.updated_at
FROM users INNER JOIN user_actions
ON users.id = user_actions.user_id
INNER JOIN actions
ON user_actions.action_id = actions.id
%s
ORDER BY user_actions.created_at DESC
LIMIT :limit;`

// Query a list of user-actions
func (r *MysqlUserActionRepo) Query(take int, filters map[string]interface{}) ([]*model.UserAction, *comtype.CommonError) {
	conditions := sqlWhereBuilder(" AND ", filters)
	filters = sqlLikeConditionFilter(filters)
	if take == 0 {
		filters["limit"] = 100
	} else {
		filters["limit"] = take
	}

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListUserAction, conditions), filters)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserActionRepo - Query:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.UserAction, 0, take)
	for rows.Next() {
		ac, err := mapUserActionRow(rows)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlUserActionRepo - Query:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, ac)
	}

	if rows.Err() != nil {
		return nil, comtype.NewCommonError(rows.Err(), "MysqlUserActionRepo - Query:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
}

func mapUserActionRow(rows *sqlx.Rows) (*model.UserAction, error) {
	var (
		ua model.UserAction
		u  model.User
		a  model.Action
	)

	err := rows.Scan(&ua.ID, &ua.UserID, &ua.ActionID, &ua.CreatedAt,
		&u.ID, &u.FullName, &u.Username, &u.Email, &u.Verified, &u.Active, &u.CreatedAt, &u.UpdatedAt,
		&a.ID, &a.ActionName, &a.ActionDesc, &a.Active, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ua.User = &u
	ua.Action = &a

	return &ua, nil
}
