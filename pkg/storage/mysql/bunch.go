package mysql

import "github.com/vespaiach/auth/pkg/adding"

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
