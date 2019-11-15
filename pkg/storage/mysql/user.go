package mysql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/modifying"
)

var sqlCreateUser = `INSERT INTO users(username, email, hash) VALUES(?, ?, ?);`

func (st *Storage) AddUser(u adding.User) (int64, error) {
	stmt, err := st.DbClient.Prepare(sqlCreateUser)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(u.Username, u.Email, u.Hash)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlCheckUsername = `SELECT count(id) FROM users WHERE username = ?;`

func (st *Storage) IsDuplicatedUsername(username string) (bool, error) {
	rows, err := st.DbClient.Queryx(sqlCheckUsername, username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	var count int
	if err := rows.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

var sqlCheckEmail = `SELECT count(id) FROM users WHERE email = ?;`

func (st *Storage) IsDuplicatedEmail(email string) (bool, error) {
	rows, err := st.DbClient.Queryx(sqlCheckEmail, email)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	var count int
	if err := rows.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

var sqlGetUsernameAndEmail = "SELECT `username`, email FROM `users` WHERE id = ?;"

func (st *Storage) GetUsernameAndEmail(id int64) (string, string, error) {
	rows, err := st.DbClient.Queryx(sqlGetUsernameAndEmail, id)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", "", errors.New(fmt.Sprintf("no user found with id = %d", id))
	}

	var username, email string
	if err := rows.Scan(&username, &email); err != nil {
		return "", "", err
	}

	return username, email, nil
}

var sqlUpdateUser = "UPDATE `users` SET %s WHERE `users`.id = :id;"

func (st *Storage) ModifyUser(u modifying.User) error {
	fields := make([]string, 0)
	updating := make(map[string]interface{})

	if len(u.Username) > 0 {
		fields = append(fields, "`username`=:username")
		updating["username"] = u.Username
	}

	if len(u.Email) > 0 {
		fields = append(fields, "`email`=:email")
		updating["email"] = u.Email
	}

	if u.Active.Valid {
		fields = append(fields, "`active`=:active")
		updating["active"] = u.Active.Bool
	}

	if len(fields) > 0 {
		fields = append(fields, "updated_at=:updated_at")
		updating["updated_at"] = time.Now()
		updating["id"] = u.ID

		_, err := st.DbClient.NamedExec(fmt.Sprintf(sqlUpdateUser, strings.Join(fields, ",")), updating)
		if err != nil {
			return err
		}
	}

	return nil
}
