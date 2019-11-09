package mysql

import "github.com/vespaiach/auth/pkg/adding"

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
