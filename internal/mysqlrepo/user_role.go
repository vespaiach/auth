package mysqlrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

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
func (r *MysqlUserRoleRepo) GetByID(id int64) (*model.UserRole, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetUserRoleByID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRoleRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.NewCommonError(nil, "MysqlUserRoleRepo - GetByID:", comtype.ErrDataNotFound, nil)
	}

	userRole, err := mapUserRoleRow(rows)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRoleRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}

	return userRole, nil
}

var sqlCreateUserRole = `
INSERT INTO user_roles(user_id, role_id) VALUES(?, ?);
`

// Create a new user-role
func (r *MysqlUserRoleRepo) Create(userID int64, roleID int64) (int64, *comtype.CommonError) {
	stmt, err := r.DbClient.Prepare(sqlCreateUserRole)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserRoleRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(userID, roleID)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserRoleRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserRoleRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return lastID, nil
}

var sqlDeleteUserRole = `
DELETE FROM user_roles WHERE user_roles.id = ?;
`

// Delete user-role
func (r *MysqlUserRoleRepo) Delete(id int64) *comtype.CommonError {
	stmt, err := r.DbClient.Prepare(sqlDeleteUserRole)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlUserRoleRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlUserRoleRepo - Delete:", comtype.ErrHandleDataFail, nil)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		return comtype.NewCommonError(err, "MysqlUserRoleRepo - Delete:", comtype.ErrHandleDataFail, nil)
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
LIMIT :limit;`

// Query a list of user-roles
func (r *MysqlUserRoleRepo) Query(take int, filters map[string]interface{}) ([]*model.UserRole, *comtype.CommonError) {
	conditions := sqlWhereBuilder(" AND ", filters)
	filters = sqlLikeConditionFilter(filters)
	if take == 0 {
		filters["limit"] = 100
	} else {
		filters["limit"] = take
	}

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListUserRole, conditions), filters)
	if err != nil {
		fmt.Println(fmt.Sprintf(sqlListUserRole, conditions))
		return nil, comtype.NewCommonError(err, "MysqlUserRoleRepo - Query:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.UserRole, 0, take)
	for rows.Next() {
		ac, err := mapUserRoleRow(rows)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlUserRoleRepo - Query:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, ac)
	}

	if rows.Err() != nil {
		return nil, comtype.NewCommonError(rows.Err(), "MysqlUserRoleRepo - Query:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
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
