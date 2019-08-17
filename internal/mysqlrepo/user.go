package mysqlrepo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
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
func (r *MysqlUserRepo) GetByID(id int64) (*model.User, error) {
	rows, err := r.DbClient.Queryx(sqlGetUserByID, id)
	if err != nil {
		log.Error("MysqlUserRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	user := new(model.User)
	err = rows.StructScan(user)
	if err != nil {
		log.Error("MysqlUserRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
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
func (r *MysqlUserRepo) GetByUsername(name string) (*model.User, error) {
	rows, err := r.DbClient.Queryx(sqlGetUserByName, name)
	if err != nil {
		log.Error("MysqlUserRepo - GetByUsername:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	user := new(model.User)
	err = rows.StructScan(user)
	if err != nil {
		log.Error("MysqlUserRepo - GetByUsername:", err)
		return nil, comtype.ErrQueryDataFailed
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
func (r *MysqlUserRepo) GetByEmail(email string) (*model.User, error) {
	rows, err := r.DbClient.Queryx(sqlGetUserByEmail, email)
	if err != nil {
		log.Error("MysqlUserRepo - GetByEmail:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	user := new(model.User)
	err = rows.StructScan(user)
	if err != nil {
		log.Error("MysqlUserRepo - GetByEmail:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return user, nil
}

var sqlCreateUser = `
INSERT INTO users(full_name, username, hashed, email) VALUES(?, ?, ?, ?);
`

// Create an new user
func (r *MysqlUserRepo) Create(fullName string, username string, hashed string, email string) (int64, error) {
	stmt, err := r.DbClient.Prepare(sqlCreateUser)
	if err != nil {
		log.Error("MysqlUserRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(fullName, username, hashed, email)
	if err != nil {
		log.Error("MysqlUserRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("MysqlUserRepo - Create:", err)
		return 0, comtype.ErrCreateDataFailed
	}

	return lastID, nil
}

var sqlUpdateUser = `
UPDATE users
%s
WHERE users.id = :id;
`

// Update user
func (r *MysqlUserRepo) Update(id int64, fields map[string]interface{}) error {
	conditions := sqlWhereBuilder(", ", fields)
	if len(conditions) == 0 {
		log.Error("MysqlUserRepo - Update:", errors.New("empty updating fields"))
		return comtype.ErrUpdateDataFailed
	}

	_, err := r.DbClient.NamedExec(fmt.Sprintf(sqlUpdateUser, conditions), fields)
	if err != nil {
		log.Error("MysqlUserRepo - Update:", err)
		return comtype.ErrUpdateDataFailed
	}

	return nil
}

const sqlListUser = `
SELECT *
FROM users
%s
ORDER BY %s 
LIMIT :offset, :limit;`

const sqlCountListUser = `
SELECT Count(*)
FROM users
%s ;`

// Query a list of users
func (r *MysqlUserRepo) Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.User, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListUser, conditions), filters)
		if err != nil {
			log.Error("MysqlUserRepo - Query:", err)
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

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListUser, conditions, sortings), filters)
	if err != nil {
		log.Error("MysqlUserRepo - Query:", err)
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.User, 0, perPage)
	for rows.Next() {
		var ac model.User
		rows.StructScan(&ac)
		if err != nil {
			log.Error("MysqlUserRepo - Query:", err)
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
