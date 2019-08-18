package mysqlrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
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
func (r *MysqlUserActionRepo) GetByID(id int64) (*model.UserAction, error) {
	rows, err := r.DbClient.Queryx(sqlGetUserActionByID, id)
	if err != nil {
		log.Error("MysqlUserActionRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	userAction, err := mapUserActionRow(rows)
	if err != nil {
		log.Error("MysqlUserActionRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return userAction, nil
}

var sqlCreateUserAction = `
INSERT INTO user_actions(user_id, action_id) VALUES(?, ?);
`

// Create a new user-action
func (r *MysqlUserActionRepo) Create(userID int64, actionID int64) (int64, error) {
	stmt, err := r.DbClient.Prepare(sqlCreateUserAction)
	if err != nil {
		log.Error("MysqlUserActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(userID, actionID)
	if err != nil {
		log.Error("MysqlUserActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlUserActionRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	return lastID, nil
}

var sqlDeleteUserAction = `
DELETE FROM user_actions WHERE user_actions.id = ?;
`

// Delete user-action
func (r *MysqlUserActionRepo) Delete(id int64) error {
	stmt, err := r.DbClient.Prepare(sqlDeleteUserAction)
	if err != nil {
		log.Error("MysqlUserActionRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	res, err := stmt.Exec(id)
	if err != nil {
		log.Error("MysqlUserActionRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		log.Error("MysqlUserActionRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
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
LIMIT :offset, :limit;`

const sqlCountListUserAction = `
SELECT Count(*)
FROM users INNER JOIN user_actions
ON users.id = user_actions.user_id
INNER JOIN actions
ON user_actions.action_id = actions.id
%s ;`

// Query a list of user-actions
func (r *MysqlUserActionRepo) Query(page int, perPage int, filters map[string]interface{}) ([]*model.UserAction, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListUserAction, conditions), filters)
		if err != nil {
			log.Error("MysqlUserActionRepo - Query:", err)
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

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListUserAction, conditions), filters)
	if err != nil {
		log.Error("MysqlUserActionRepo - Query:", err)
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.UserAction, 0, perPage)
	for rows.Next() {
		ac, err := mapUserActionRow(rows)
		if err != nil {
			log.Error("MysqlUserActionRepo - Query:", err)
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
