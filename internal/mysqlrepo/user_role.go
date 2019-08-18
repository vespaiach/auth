package mysqlrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlUserRoleRepo will implement model.UserRoleRepo
type MysqlUserRoleRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlUserRoleRepo create new instance of MysqlUserRoleRepo
func NewMysqlUserRoleRepo(db *sqlx.DB) model.UserRoleRepo {
	return &MysqlUserRoleRepo{
		db,
	}
}

var sqlGetUserRoleByID = `
SELECT
	user_roles.id,
	user_roles.user_id,
	user_roles.role_id,
	user_roles.created_at,
	users.id,
	users.full_name,
	users.username,
	users.email,
	users.verified,
	users.active,
	users.created_at,
	users.updated_at,
	roles.id,
	roles.role_name,
	roles.role_desc,
	roles.active,
	roles.created_at,
	roles.updated_at
FROM users INNER JOIN user_roles
ON users.id = user_roles.user_id
INNER JOIN roles
ON user_roles.role_id = roles.id
WHERE user_roles.id =?
LIMIT 1;
`

// GetByID find a user-role by its ID
func (r *MysqlUserRoleRepo) GetByID(id int64) (*model.UserRole, error) {
	rows, err := r.DbClient.Queryx(sqlGetUserRoleByID, id)
	if err != nil {
		log.Error("MysqlUserRoleRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	userRole, err := mapUserRoleRow(rows)
	if err != nil {
		log.Error("MysqlUserRoleRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return userRole, nil
}

var sqlCreateUserRole = `
INSERT INTO user_roles(user_id, role_id) VALUES(?, ?);
`

// Create a new user-role
func (r *MysqlUserRoleRepo) Create(userID int64, roleID int64) (int64, error) {
	stmt, err := r.DbClient.Prepare(sqlCreateUserRole)
	if err != nil {
		log.Error("MysqlUserRoleRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(userID, roleID)
	if err != nil {
		log.Error("MysqlUserRoleRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlUserRoleRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	return lastID, nil
}

var sqlDeleteUserRole = `
DELETE FROM user_roles WHERE user_roles.id = ?;
`

// Delete user-role
func (r *MysqlUserRoleRepo) Delete(id int64) error {
	stmt, err := r.DbClient.Prepare(sqlDeleteUserRole)
	if err != nil {
		log.Error("MysqlUserRoleRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	res, err := stmt.Exec(id)
	if err != nil {
		log.Error("MysqlUserRoleRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		log.Error("MysqlUserRoleRepo - Delete:", err)
		return comtype.ErrDeleteDataFailed
	}

	return nil
}

const sqlListUserRole = `
SELECT
	user_roles.id,
	user_roles.user_id,
	user_roles.role_id,
	user_roles.created_at,
	users.id,
	users.full_name,
	users.username,
	users.email,
	users.verified,
	users.active,
	users.created_at,
	users.updated_at,
	roles.id,
	roles.role_name,
	roles.role_desc,
	roles.active,
	roles.created_at,
	roles.updated_at
FROM users INNER JOIN user_roles
ON users.id = user_roles.user_id
INNER JOIN roles
ON user_roles.role_id = roles.id
%s
ORDER BY user_roles.created_at DESC
LIMIT :offset, :limit;`

const sqlCountListUserRole = `
SELECT Count(*)
FROM users INNER JOIN user_roles
ON users.id = user_roles.user_id
INNER JOIN roles
ON user_roles.role_id = roles.id
%s ;`

// Query a list of user-roles
func (r *MysqlUserRoleRepo) Query(page int, perPage int, filters map[string]interface{}) ([]*model.UserRole, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListUserRole, conditions), filters)
		if err != nil {
			log.Error("MysqlUserRoleRepo - Query:", err)
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

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListUserRole, conditions), filters)
	if err != nil {
		log.Error("MysqlUserRoleRepo - Query:", err)
		fmt.Println(fmt.Sprintf(sqlListUserRole, conditions))
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.UserRole, 0, perPage)
	for rows.Next() {
		ac, err := mapUserRoleRow(rows)
		if err != nil {
			log.Error("MysqlUserRoleRepo - Query:", err)
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

func mapUserRoleRow(rows *sqlx.Rows) (*model.UserRole, error) {
	var (
		ur model.UserRole
		u  model.User
		r  model.Role
	)

	err := rows.Scan(&ur.ID, &ur.UserID, &ur.RoleID, &ur.CreatedAt,
		&u.ID, &u.FullName, &u.Username, &u.Email, &u.Verified, &u.Active, &u.CreatedAt, &u.UpdatedAt,
		&r.ID, &r.RoleName, &r.RoleDesc, &r.Active, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ur.User = &u
	ur.Role = &r

	return &ur, nil
}
