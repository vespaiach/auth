package mysql

import (
	"errors"
	"fmt"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/modifying"
	"strings"
	"time"
)

var sqlCreateBunch = "INSERT INTO `bunch` (`name`, `desc`) VALUES (?, ?);"

func (st *Storage) AddBunch(b adding.Bunch) (int64, error) {
	stmt, err := st.DbClient.Prepare(sqlCreateBunch)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(b.Name, b.Desc)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlCheckBunchName = `SELECT count(id) FROM bunch WHERE name = ?;`

func (st *Storage) IsDuplicatedBunch(name string) (bool, error) {
	rows, err := st.DbClient.Queryx(sqlCheckBunchName, name)
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

var sqlGetBunchName = `SELECT name FROM bunch WHERE id = ?;`

func (st *Storage) GetBunchNameByID(id int64) (string, error) {
	rows, err := st.DbClient.Queryx(sqlGetBunchName, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New(fmt.Sprintf("no bunch found with id = %d", id))
	}

	var name string
	if err := rows.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

var sqlUpdateBunch = "UPDATE bunch SET %s WHERE bunch.id = :id;"

func (st *Storage) ModifyBunch(b modifying.Bunch) error {
	fields := make([]string, 0)
	updating := make(map[string]interface{})

	if len(b.Name) > 0 {
		fields = append(fields, "`name`=:name")
		updating["name"] = b.Name
	}

	if len(b.Desc) > 0 {
		fields = append(fields, "`desc`=:desc")
		updating["desc"] = b.Desc
	}

	if b.Active.Valid {
		fields = append(fields, "`active`=:active")
		updating["active"] = b.Active.Bool
	}

	if len(fields) > 0 {
		fields = append(fields, "`updated_at`=:updated_at")
		updating["updated_at"] = time.Now()
		updating["id"] = b.ID

		_, err := st.DbClient.NamedExec(fmt.Sprintf(sqlUpdateBunch, strings.Join(fields, ",")), updating)
		if err != nil {
			return err
		}
	}

	return nil
}
