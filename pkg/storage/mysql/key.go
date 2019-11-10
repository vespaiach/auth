package mysql

import (
	"errors"
	"fmt"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/modifying"
	"strings"
	"time"
)

var sqlCreateKey = "INSERT INTO `keys` (`key`, `desc`) VALUES (?, ?);"

func (st *Storage) AddServiceKey(sk adding.ServiceKey) (int64, error) {
	stmt, err := st.DbClient.Prepare(sqlCreateKey)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(sk.Key, sk.Desc)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlCheckKey = `SELECT count(id) FROM keys WHERE key = ?;`

func (st *Storage) IsDuplicatedKey(key string) (bool, error) {
	rows, err := st.DbClient.Queryx(sqlCheckKey, key)
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

var sqlGetKeyName = "SELECT `key` FROM `keys` WHERE id = ?;"

func (st *Storage) GetKeyByID(id int64) (string, error) {
	rows, err := st.DbClient.Queryx(sqlGetKeyName, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New(fmt.Sprintf("no key found with id = %d", id))
	}

	var key string
	if err := rows.Scan(&key); err != nil {
		return "", err
	}

	return key, nil
}

var sqlUpdateKeys = "UPDATE `keys` SET %s WHERE `keys`.id = :id;"

func (st *Storage) ModifyServiceKey(sk modifying.ServiceKey) error {
	fields := make([]string, 0)
	updating := make(map[string]interface{})

	if len(sk.Key) > 0 {
		fields = append(fields, "`key`=:key")
		updating["key"] = sk.Key
	}

	if len(sk.Desc) > 0 {
		fields = append(fields, "`desc`=:desc")
		updating["desc"] = sk.Desc
	}

	if len(fields) > 0 {
		fields = append(fields, "`updated_at`=:updated_at")
		updating["updated_at"] = time.Now()
		updating["id"] = sk.ID

		_, err := st.DbClient.NamedExec(fmt.Sprintf(sqlUpdateKeys, strings.Join(fields, ",")), updating)
		if err != nil {
			return err
		}
	}

	return nil
}
