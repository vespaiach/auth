package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlUserRepo will implement model.UserRepo
type MysqlUserRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlUserRepo create new instance of MysqlUserRepo
func NewMysqlUserRepo(db *sqlx.DB) model.UserRepo {
	return &MysqlUserRepo{
		db,
	}
}

var sqlGetUserByID = `
SELECT *
FROM users
WHERE users.id =?
LIMIT 1;
`

// GetByID find an user by its ID
func (r *MysqlUserRepo) GetByID(id int64) (*model.User, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetUserByID, id)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.NewCommonError(nil, "MysqlUserRepo - GetByID:", comtype.ErrDataNotFound, nil)
	}

	user := new(model.User)
	err = rows.StructScan(user)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - GetByID:", comtype.ErrQueryDataFail, nil)
	}

	return user, nil
}

var sqlGetUserByName = `
SELECT *
FROM users
WHERE users.username = ?
LIMIT 1;
`

// GetByUsername find an user by its name
func (r *MysqlUserRepo) GetByUsername(name string) (*model.User, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetUserByName, name)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - GetByUsername:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	user := new(model.User)
	err = rows.StructScan(user)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - GetByUsername:", comtype.ErrQueryDataFail, nil)
	}

	return user, nil
}

var sqlGetUserByEmail = `
SELECT *
FROM users
WHERE users.email = ?
LIMIT 1;
`

// GetByEmail find an user by its email
func (r *MysqlUserRepo) GetByEmail(email string) (*model.User, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetUserByEmail, email)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - GetByEmail:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	user := new(model.User)
	err = rows.StructScan(user)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - GetByEmail:", comtype.ErrQueryDataFail, nil)
	}

	return user, nil
}

var sqlCreateUser = `
INSERT INTO users(full_name, username, hashed, email) VALUES(?, ?, ?, ?);
`

// Create an new user
func (r *MysqlUserRepo) Create(fullName string, username string, hashed string, email string) (int64, *comtype.CommonError) {
	stmt, err := r.DbClient.Prepare(sqlCreateUser)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(fullName, username, hashed, email)
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, comtype.NewCommonError(err, "MysqlUserRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return lastID, nil
}

var sqlUpdateUser = `
UPDATE users SET 
%s
WHERE users.id = :id;
`

// Update user
func (r *MysqlUserRepo) Update(id int64, fields map[string]interface{}) *comtype.CommonError {
	if len(fields) == 0 {
		return comtype.NewCommonError(errors.New("empty updating fields"), "MysqlUserRepo - Update:",
			comtype.ErrHandleDataFail, nil)
	}

	fields["id"] = id
	_, err := r.DbClient.NamedExec(fmt.Sprintf(sqlUpdateUser, sqlUpdateBuilder(fields, map[string]bool{"id": true})), fields)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlUserRepo - Update:", comtype.ErrHandleDataFail, nil)
	}

	return nil
}

const sqlListUser = `
SELECT *
FROM users
%s
ORDER BY %s 
LIMIT :limit;`

// Query a list of users
func (r *MysqlUserRepo) Query(take int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.User, *comtype.CommonError) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	if take == 0 {
		filters["limit"] = 100
	} else {
		filters["limit"] = take
	}

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListUser, conditions, sortings), filters)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlUserRepo - Query:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.User, 0, take)
	for rows.Next() {
		var ac model.User
		rows.StructScan(&ac)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlUserRepo - Query:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, &ac)
	}

	if rows.Err() != nil {
		return nil, comtype.NewCommonError(rows.Err(), "MysqlUserRepo - Query:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
}
