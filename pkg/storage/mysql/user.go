package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/usrmgr"
	"strings"
	"sync"
	"time"
)

// UserStorage implements db's storage for user
type UserStorage struct {
	db *sqlx.DB
}

// NewUserStorage create new instance of BunchStorage
func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db,
	}
}

var sqlAddUser = "INSERT INTO users(username, email, hash) VALUES(?, ?, ?);"

func (st *UserStorage) AddUser(username string, email string, hash string) (int64, error) {
	stmt, err := st.db.Prepare(sqlAddUser)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(username, email, hash)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlUpdateUser = "UPDATE `users` SET %s	WHERE id = :id;"

func (st *UserStorage) ModifyUser(id int64, username string, email string, hash string, active sql.NullBool) error {
	updating := make(map[string]interface{})
	var condition string
	var prefix string

	if len(username) > 0 {
		updating["username"] = username
		condition += prefix + "`username` = :username"
		prefix = ", "
	}

	if len(email) > 0 {
		updating["email"] = email
		condition += prefix + "`email` = :email"
		prefix = ", "
	}

	if len(hash) > 0 {
		updating["hash"] = hash
		condition += prefix + "`hash` = :hash"
		prefix = ", "
	}

	if active.Valid {
		updating["active"] = active.Bool
		condition += prefix + "`active` = :active"
		prefix = ", "
	}

	if len(updating) > 0 {
		updating["id"] = id
		updating["updated_at"] = time.Now()
		condition += prefix + "`updated_at` = :updated_at"

		_, err := st.db.NamedExec(fmt.Sprintf(sqlUpdateUser, condition), updating)
		if err != nil {
			return err
		}
	}

	return nil
}

var sqlGetUserByName = "SELECT id, `username`, `email`, `hash`, active, created_at, updated_at FROM `users` " +
	"WHERE `username` = ? LIMIT 1;"

func (st *UserStorage) GetUserByUsername(username string) (*usrmgr.User, error) {
	rows, err := st.db.Queryx(sqlGetUserByName, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := new(usrmgr.User)
	err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Hash, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

var sqlGetUserEmail = "SELECT id, `username`, `email`, `hash`, active, created_at, updated_at FROM `users` " +
	"WHERE `email` = ? LIMIT 1;"

func (st *UserStorage) GetUserByEmail(email string) (*usrmgr.User, error) {
	rows, err := st.db.Queryx(sqlGetUserEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := new(usrmgr.User)
	err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Hash, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

var sqlGetUserByID = "SELECT id, `username`, `email`, `hash`, active, created_at, updated_at FROM `users` " +
	"WHERE `id` = ? LIMIT 1;"

func (st *UserStorage) GetUser(id int64) (*usrmgr.User, error) {
	rows, err := st.db.Queryx(sqlGetUserByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := new(usrmgr.User)
	err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Hash, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

var sqlAddBunchesToUser = "INSERT INTO `user_bunches` (user_id, bunch_id, created_at) VALUES %s;"

func (st *UserStorage) AddBunchesToUser(userID int64, bunchIDs []int64) error {
	updating := make([]string, 0, len(bunchIDs))
	for _, id := range bunchIDs {
		updating = append(updating, fmt.Sprintf("(%d, %d, :created_at)", userID, id))
	}

	sql := fmt.Sprintf(sqlAddBunchesToUser, strings.Join(updating, ", "))

	_, err := st.db.NamedExec(sql, map[string]interface{}{"created_at": time.Now()})
	if err != nil {
		return err
	}

	return nil
}

var sqlQueryUsers = "SELECT id, `username`, `email`, `hash`, active, created_at, updated_at FROM `users` %s ORDER BY %s LIMIT :offset, :limit;"
var sqlQueryUsersCounter = "SELECT count(id) FROM `users` %s;"

func (st *UserStorage) QueryUsers(take int64, skip int64, username string, email string, active sql.NullBool,
	sortby string, direction common.SortingDirection) ([]*usrmgr.User, int64, error) {

	var (
		order         string
		where         string
		sql           string
		filter        map[string]interface{}
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*usrmgr.User
		total         int64
		prefix        string
	)

	filter = make(map[string]interface{})

	if len(username) > 0 {
		where = "WHERE `username` LIKE :username"
		filter["username"] = "%" + username + "%"
		prefix = " AND "
	} else {
		prefix = " WHERE "
	}

	if len(email) > 0 {
		where += prefix + "`email` LIKE :email"
		filter["email"] = "%" + email + "%"
		prefix = " AND "
	}

	if active.Valid {
		where += prefix + "`active` = :active"
		filter["active"] = active.Bool
	}

	if direction == common.Descending {
		order = fmt.Sprintf("`%s` DESC, id", sortby)
	} else {
		order = fmt.Sprintf("`%s` ASC, id", sortby)
	}

	sql = fmt.Sprintf(sqlQueryUsers, where, order)

	filter["offset"] = skip
	filter["limit"] = take

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := st.db.NamedQuery(sql, filter)
		if err != nil {
			queryErr = err
			return
		}
		defer rows.Close()

		results = make([]*usrmgr.User, 0, take)
		for rows.Next() {
			u := new(usrmgr.User)
			err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Hash, &u.Active, &u.CreatedAt, &u.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, u)
		}

		if rows.Err() != nil {
			queryErr = rows.Err()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := st.db.NamedQuery(fmt.Sprintf(sqlQueryUsersCounter, where), filter)
		if err != nil {
			countTotalErr = err
			return
		}
		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&total)
			if err != nil {
				countTotalErr = err
				return
			}
		}
	}()

	wg.Wait()

	if queryErr != nil {
		return nil, 0, queryErr
	}
	if countTotalErr != nil {
		return nil, 0, countTotalErr
	}

	return results, total, nil
}

var sqlGetBunchIDs = "SELECT id FROM bunches WHERE `name` IN (%s);"

func (st *UserStorage) GetBunchIDs(bunches []string) ([]int64, error) {
	len := len(bunches)
	conditions := make([]string, 0, len)
	values := make([]interface{}, 0, len)
	for i := 0; i < len; i++ {
		conditions = append(conditions, "?")
		values = append(values, interface{}(bunches[i]))
	}

	sql := fmt.Sprintf(sqlGetBunchIDs, strings.Join(conditions, ","))
	rows, err := st.db.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		results = append(results, id)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}

var sqlGetBunchesByUsername = "SELECT bunches.id, bunches.`name`, bunches.`desc`, bunches.active, " +
	"bunches.created_at, bunches.updated_at FROM `user_bunches` " +
	"INNER JOIN `users` ON `user_bunches`.user_id = `users`.id " +
	"INNER JOIN `bunches` ON `bunches`.id = `user_bunches`.bunch_id " +
	"WHERE `users`.username = ?"

func (st *UserStorage) GetBunches(username string) ([]*usrmgr.Bunch, error) {
	rows, err := st.db.Queryx(sqlGetBunchesByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*usrmgr.Bunch, 0)
	for rows.Next() {
		b := new(usrmgr.Bunch)
		if err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}

var sqlGetKeysByUsername = "SELECT DISTINCT `keys`.id, `keys`.`key`, `keys`.`desc`, `keys`.created_at, `keys`.updated_at " +
	"FROM `users` INNER JOIN user_bunches ON `user_bunches`.user_id = `users`.id " +
	"INNER JOIN bunches ON bunches.id = user_bunches.bunch_id " +
	"INNER JOIN bunch_keys ON bunch_keys.bunch_id = bunches.id " +
	"INNER JOIN `keys` ON `keys`.id = bunch_keys.key_id " +
	"WHERE `users`.username = ?"

func (st *UserStorage) GetKeys(username string) ([]*usrmgr.Key, error) {
	rows, err := st.db.Queryx(sqlGetKeysByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*usrmgr.Key, 0)
	for rows.Next() {
		key := new(usrmgr.Key)
		if err := rows.Scan(&key.ID, &key.Key, &key.Desc, &key.CreatedAt, &key.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}
