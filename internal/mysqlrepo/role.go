package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
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
func (r *MysqlRoleRepo) GetByID(id int64) (*model.Role, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetRoleByID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	role := new(model.Role)
	err = rows.StructScan(role)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}

	return role, nil
}

var sqlGetRoleByUserID = `
SELECT 
	roles.id,
	roles.role_name,
	roles.role_desc,
	roles.active,
	roles.created_at,
	roles.updated_at
FROM roles INNER JOIN user_roles
ON roles.id = user_roles.role_id
WHERE user_roles.user_id = ? AND active = 1;
`

// GetByUserID gets role by user ID
func (r *MysqlRoleRepo) GetByUserID(id int64) ([]*model.Role, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetRoleByUserID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.Role, 0)
	for rows.Next() {
		var ac model.Role
		rows.StructScan(&ac)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlRoleRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, &ac)
	}

	return results, nil
}

var sqlGetRoleByName = `
SELECT *
FROM roles
WHERE roles.role_name = ?
LIMIT 1;
`

// GetByName find a role by its name
func (r *MysqlRoleRepo) GetByName(name string) (*model.Role, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetRoleByName, name)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleRepo - GetByName:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.NewCommonError(nil, "MysqlRoleRepo - GetByName:", comtype.ErrDataNotFound, nil)
	}

	role := new(model.Role)
	err = rows.StructScan(role)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleRepo - GetByName:", comtype.ErrQueryDataFail, nil)
	}

	return role, nil
}

var sqlCreateRole = `
INSERT INTO roles(role_name, role_desc) VALUES(?, ?);
`

// Create a new role
func (r *MysqlRoleRepo) Create(name string, desc string) (int64, *comtype.CommonError) {
	stmt, err := r.DbClient.Prepare(sqlCreateRole)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlRoleRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(name, desc)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlRoleRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlRoleRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return lastID, nil
}

var sqlUpdateRole = `
UPDATE roles SET 
%s
WHERE roles.id = :id;
`

// Update role
func (r *MysqlRoleRepo) Update(id int64, fields map[string]interface{}) *comtype.CommonError {
	if len(fields) == 0 {
		return comtype.NewCommonError(errors.New("empty updating fields"), "MysqlRoleRepo - Update:",
			comtype.ErrHandleDataFail, nil)
	}

	fields["id"] = id
	_, err := r.DbClient.NamedExec(fmt.Sprintf(sqlUpdateRole, sqlUpdateBuilder(fields, map[string]bool{"id": true})), fields)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlRoleRepo - Update:", comtype.ErrHandleDataFail, nil)
	}

	return nil
}

const sqlListRole = `
SELECT *
FROM roles
%s
ORDER BY %s 
LIMIT :limit;`

// Query a list of roles
func (r *MysqlRoleRepo) Query(take int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.Role, *comtype.CommonError) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	if take == 0 {
		filters["limit"] = 100
	} else {
		filters["limit"] = take
	}

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListRole, conditions, sortings), filters)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlRoleRepo - Query:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.Role, 0, take)
	for rows.Next() {
		var ac model.Role
		rows.StructScan(&ac)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlRoleRepo - Query:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, &ac)
	}

	if rows.Err() != nil {
		return nil, comtype.NewCommonError(rows.Err(), "MysqlRoleRepo - Query:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
}
