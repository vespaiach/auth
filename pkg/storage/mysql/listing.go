package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/pkg/listing"
)

type ListingDbHelper struct {
	DbClient *sqlx.DB
}

var sqlGetUserByID = "SELECT `id`, `username`, `email`, active, hash, created_at, updated_at FROM `users` " +
	"WHERE id = ? LIMIT 1;"

func (st *ListingDbHelper) GetUserByID(id int64) (*listing.User, error) {
	rows, err := st.DbClient.Queryx(sqlGetUserByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := new(listing.User)
	err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Active, &u.Hash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

var sqlGetUserByUsername = "SELECT `id`, `username`, `email`, active, hash, created_at, updated_at FROM `users` " +
	"WHERE username = ? LIMIT 1;"

func (st *ListingDbHelper) GetUserByUsername(username string) (*listing.User, error) {
	rows, err := st.DbClient.Queryx(sqlGetUserByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := new(listing.User)
	err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Active, &u.Hash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

var sqlGetUserByEmail = "SELECT `id`, `username`, `email`, active, hash, created_at, updated_at FROM `users` " +
	"WHERE email = ? LIMIT 1;"

func (st *ListingDbHelper) GetUserByEmail(email string) (*listing.User, error) {
	rows, err := st.DbClient.Queryx(sqlGetUserByEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := new(listing.User)
	err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Active, &u.Hash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

var sqlGetBunchByID = "SELECT `id`, `name`, `desc`, active, created_at, updated_at FROM `bunch` " +
	"WHERE id = ? LIMIT 1;"

func (st *ListingDbHelper) GetBunchByID(id int64) (*listing.Bunch, error) {
	rows, err := st.DbClient.Queryx(sqlGetBunchByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := new(listing.Bunch)
	err = rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var sqlGetBunchByName = "SELECT `id`, `name`, `desc`, active, created_at, updated_at FROM `bunch` " +
	"WHERE `name` = ? LIMIT 1;"

func (st *ListingDbHelper) GetBunchByName(name string) (*listing.Bunch, error) {
	rows, err := st.DbClient.Queryx(sqlGetBunchByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := new(listing.Bunch)
	err = rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var sqlGetKeyByID = "SELECT `id`, `key`, `desc`, created_at, updated_at FROM `keys` " +
	"WHERE id = ? LIMIT 1;"

func (st *ListingDbHelper) GetKeyByID(id int64) (*listing.Key, error) {
	rows, err := st.DbClient.Queryx(sqlGetKeyByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := new(listing.Key)
	err = rows.Scan(&b.ID, &b.Key, &b.Desc, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var sqlGetKeyByName = "SELECT `id`, `key`, `desc`, created_at, updated_at FROM `keys` " +
	"WHERE `key` = ? LIMIT 1;"

func (st *ListingDbHelper) GetKeyByName(key string) (*listing.Key, error) {
	rows, err := st.DbClient.Queryx(sqlGetKeyByName, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := new(listing.Key)
	err = rows.Scan(&b.ID, &b.Key, &b.Desc, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// var sqlCreateBunch = "INSERT INTO `bunch` (`name`, `desc`) VALUES (?, ?);"

// func (st *ListingDbHelper) AddBunch(b adding.Bunch) (int64, error) {
// 	stmt, err := st.DbClient.Prepare(sqlCreateBunch)
// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := stmt.Exec(b.Name, b.Desc)
// 	if err != nil {
// 		return 0, err
// 	}

// 	lastID, err := res.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return lastID, nil
// }

// var sqlCheckBunchName = `SELECT count(id) FROM bunch WHERE name = ?;`

// func (st *ListingDbHelper) IsDuplicatedBunch(name string) (bool, error) {
// 	rows, err := st.DbClient.Queryx(sqlCheckBunchName, name)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer rows.Close()

// 	if !rows.Next() {
// 		return false, nil
// 	}

// 	var count int
// 	if err := rows.Scan(&count); err != nil {
// 		return false, err
// 	}

// 	return count > 0, nil
// }

// var sqlGetBunchName = `SELECT name FROM bunch WHERE id = ?;`

// func (st *ListingDbHelper) GetBunchNameByID(id int64) (string, error) {
// 	rows, err := st.DbClient.Queryx(sqlGetBunchName, id)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer rows.Close()

// 	if !rows.Next() {
// 		return "", errors.New(fmt.Sprintf("no bunch found with id = %d", id))
// 	}

// 	var name string
// 	if err := rows.Scan(&name); err != nil {
// 		return "", err
// 	}

// 	return name, nil
// }

// var sqlUpdateBunch = "UPDATE bunch SET %s WHERE bunch.id = :id;"

// func (st *ListingDbHelper) ModifyBunch(b modifying.Bunch) error {
// 	fields := make([]string, 0)
// 	updating := make(map[string]interface{})

// 	if len(b.Name) > 0 {
// 		fields = append(fields, "`name`=:name")
// 		updating["name"] = b.Name
// 	}

// 	if len(b.Desc) > 0 {
// 		fields = append(fields, "`desc`=:desc")
// 		updating["desc"] = b.Desc
// 	}

// 	if b.Active.Valid {
// 		fields = append(fields, "`active`=:active")
// 		updating["active"] = b.Active.Bool
// 	}

// 	if len(fields) > 0 {
// 		fields = append(fields, "`updated_at`=:updated_at")
// 		updating["updated_at"] = time.Now()
// 		updating["id"] = b.ID

// 		_, err := st.DbClient.NamedExec(fmt.Sprintf(sqlUpdateBunch, strings.Join(fields, ",")), updating)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
