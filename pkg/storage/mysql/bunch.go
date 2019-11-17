package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/pkg/bunchmgr"
	"github.com/vespaiach/auth/pkg/common"
	"strings"
	"sync"
	"time"
)

// BunchStorage implements db's storage for bunch
type BunchStorage struct {
	db *sqlx.DB
}

// NewBunchStorage create new instance of BunchStorage
func NewBunchStorage(db *sqlx.DB) *BunchStorage {
	return &BunchStorage{
		db,
	}
}

var sqlCreateBunch = "INSERT INTO bunches (`name`, `desc`, `active`, created_at, updated_at) VALUES (?, ?, ?, ?, ?);"

func (st *BunchStorage) AddBunch(name string, desc string) (int64, error) {
	stmt, err := st.db.Prepare(sqlCreateBunch)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	res, err := stmt.Exec(name, desc, true, now, now)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlGetBunchByName = "SELECT id, `name`, `desc`, active, created_at, updated_at FROM `bunches` WHERE `name` = ? LIMIT 1;"

func (st *BunchStorage) GetBunchByName(name string) (*bunchmgr.Bunch, error) {
	rows, err := st.db.Queryx(sqlGetBunchByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := new(bunchmgr.Bunch)
	if err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, err
	}

	return b, nil
}

var sqlGetBunchByID = "SELECT id, `name`, `desc`, active, created_at, updated_at FROM `bunches` WHERE id = ? LIMIT 1;"

func (st *BunchStorage) GetBunch(id int64) (*bunchmgr.Bunch, error) {
	rows, err := st.db.Queryx(sqlGetBunchByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := new(bunchmgr.Bunch)
	if err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, err
	}

	return b, nil
}

var sqlUpdateBunch = "UPDATE bunches SET %s	WHERE id = :id;"

func (st *BunchStorage) ModifyBunch(id int64, name string, desc string, active sql.NullBool) error {
	updating := make(map[string]interface{})
	var condition string
	var prefix string

	if len(name) > 0 {
		updating["name"] = name
		condition += prefix + "`name` = :name"
		prefix = ", "
	}

	if len(desc) > 0 {
		updating["desc"] = desc
		condition += prefix + "`desc` = :desc"
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

		_, err := st.db.NamedExec(fmt.Sprintf(sqlUpdateBunch, condition), updating)
		if err != nil {
			return err
		}
	}

	return nil
}

var sqlQueryBunches = "SELECT id, `name`, `desc`, active, created_at, updated_at FROM `bunches` %s ORDER BY %s LIMIT :offset, :limit;"
var sqlQueryBunchesCounter = "SELECT count(id) FROM `bunches` %s;"

func (st *BunchStorage) QueryBunches(take int64, skip int64, name string, active sql.NullBool, sortby string,
	direction common.SortingDirection) ([]*bunchmgr.Bunch, int64, error) {

	var (
		order         string
		where         string
		sql           string
		filter        map[string]interface{}
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*bunchmgr.Bunch
		total         int64
		prefix        string
	)

	filter = make(map[string]interface{})

	if len(name) > 0 {
		where = "WHERE `name` LIKE :name"
		filter["name"] = "%" + name + "%"
		prefix = " AND "
	} else {
		prefix = " WHERE "
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

	sql = fmt.Sprintf(sqlQueryBunches, where, order)

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

		results = make([]*bunchmgr.Bunch, 0, take)
		for rows.Next() {
			b := new(bunchmgr.Bunch)
			err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active, &b.CreatedAt, &b.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, b)
		}

		if rows.Err() != nil {
			queryErr = rows.Err()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := st.db.NamedQuery(fmt.Sprintf(sqlQueryBunchesCounter, where), filter)
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

var sqlAddKeysToBunch = "INSERT INTO `bunch_keys` (bunch_id, key_id, created_at) VALUES %s;"

func (st *BunchStorage) AddKeysToBunch(bunchID int64, keyIDs []int64) error {
	updating := make([]string, 0, len(keyIDs))
	for _, id := range keyIDs {
		updating = append(updating, fmt.Sprintf("(%d, %d, :created_at)", bunchID, id))
	}

	sql := fmt.Sprintf(sqlAddKeysToBunch, strings.Join(updating, ", "))

	_, err := st.db.NamedExec(sql, map[string]interface{}{"created_at": time.Now()})
	if err != nil {
		return err
	}

	return nil
}

var sqlGetKeyInBunch = "SELECT `keys`.id, `keys`.`key`, `keys`.`desc`, `keys`.created_at, `keys`.updated_at " +
	"FROM bunch_keys " +
	"INNER JOIN `keys` ON `keys`.id = bunch_keys.key_id " +
	"INNER JOIN `bunches` ON `bunches`.id = bunch_keys.bunch_id " +
	"WHERE bunches.name = ?"

func (st *BunchStorage) GetKeysInBunch(name string) ([]*bunchmgr.Key, error) {
	rows, err := st.db.Queryx(sqlGetKeyInBunch, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*bunchmgr.Key, 0)
	for rows.Next() {
		key := new(bunchmgr.Key)
		err := rows.Scan(&key.ID, &key.Key, &key.Desc, &key.CreatedAt, &key.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, key)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}

var sqlGetKeyIDsByKeyName = "SELECT id FROM `keys` WHERE `key` IN (%s)"

func (st *BunchStorage) GetKeyIDs(keys []string) ([]int64, error) {
	len := len(keys)
	conditions := make([]string, 0, len)
	values := make([]interface{}, 0, len)
	for i := 0; i < len; i++ {
		conditions = append(conditions, "?")
		values = append(values, interface{}(keys[i]))
	}

	sql := fmt.Sprintf(sqlGetKeyIDsByKeyName, strings.Join(conditions, ","))
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
