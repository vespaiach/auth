package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/keymgr"
	"sync"
	"time"
)

// KeyStorage implements db's storage for key
type KeyStorage struct {
	db *sqlx.DB
}

// NewKeyStorage create new instance of KeyStorage
func NewKeyStorage(db *sqlx.DB) *KeyStorage {
	return &KeyStorage{
		db,
	}
}

var sqlCreateKey = "INSERT INTO `keys` (`key`, `desc`, created_at, updated_at) VALUES (?, ?, ?, ?);"

func (st *KeyStorage) AddKey(name string, desc string) (int64, error) {
	stmt, err := st.db.Prepare(sqlCreateKey)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	res, err := stmt.Exec(name, desc, now, now)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlGetKeyByName = "SELECT id, `key`, `desc`, created_at, updated_at FROM `keys` WHERE `key` = ? LIMIT 1;"

func (st *KeyStorage) GetKeyByName(name string) (*keymgr.Key, error) {
	rows, err := st.db.Queryx(sqlGetKeyByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	key := new(keymgr.Key)
	if err := rows.Scan(&key.ID, &key.Key, &key.Desc, &key.CreatedAt, &key.UpdatedAt); err != nil {
		return nil, err
	}

	return key, nil
}

var sqlGetKeyByID = "SELECT id, `key`, `desc`, created_at, updated_at FROM `keys` WHERE id = ? LIMIT 1;"

func (st *KeyStorage) GetKey(id int64) (*keymgr.Key, error) {
	rows, err := st.db.Queryx(sqlGetKeyByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	key := new(keymgr.Key)
	if err := rows.Scan(&key.ID, &key.Key, &key.Desc, &key.CreatedAt, &key.UpdatedAt); err != nil {
		return nil, err
	}

	return key, nil
}

var sqlGetBunchID = "SELECT id FROM `bunches` WHERE `name` = ? LIMIT 1;"

func (st *KeyStorage) GetBunchID(name string) (int64, error) {
	rows, err := st.db.Queryx(sqlGetBunchID, name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}

	var id int64
	if err := rows.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

var sqlUpdateKeys = "UPDATE `keys` SET `key` = :key, `desc` = :desc, updated_at = :updated_at WHERE `keys`.id = :id;"

func (st *KeyStorage) ModifyKey(id int64, name string, desc string) error {
	updating := map[string]interface{}{
		"key":        name,
		"desc":       desc,
		"id":         id,
		"updated_at": time.Now(),
	}

	_, err := st.db.NamedExec(sqlUpdateKeys, updating)
	if err != nil {
		return err
	}

	return nil
}

var sqlAddKeyToBunch = "INSERT INTO `bunch_keys` (bunch_id, key_id, created_at) VALUES (:bunch_id, :key_id, :created_at);"

func (st *KeyStorage) AddKeyToBunch(keyID int64, bunchID int64) (int64, error) {
	updating := map[string]interface{}{
		"bunch_id":   bunchID,
		"key_id":     keyID,
		"created_at": time.Now(),
	}

	res, err := st.db.NamedExec(sqlAddKeyToBunch, updating)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

var sqlQueryKeys = "SELECT id, `key`, `desc`, created_at, updated_at FROM `keys` %s ORDER BY %s LIMIT :offset, :limit;"
var sqlQueryKeysCounter = "SELECT count(id) FROM `keys` %s;"

func (st *KeyStorage) QueryKeys(take int64, skip int64, name string, sortby string,
	direction common.SortingDirection) ([]*keymgr.Key, int64, error) {

	var (
		order         string
		where         string
		sql           string
		filter        map[string]interface{}
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*keymgr.Key
		total         int64
	)

	filter = make(map[string]interface{})

	if len(name) > 0 {
		where = "WHERE `key` LIKE :name"
		filter["name"] = "%" + name + "%"
	}

	if direction == common.Descending {
		order = fmt.Sprintf("`%s` DESC, id", sortby)
	} else {
		order = fmt.Sprintf("`%s` ASC, id", sortby)
	}

	sql = fmt.Sprintf(sqlQueryKeys, where, order)

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

		results = make([]*keymgr.Key, 0, take)
		for rows.Next() {
			key := new(keymgr.Key)
			err := rows.Scan(&key.ID, &key.Key, &key.Desc, &key.CreatedAt, &key.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, key)
		}

		if rows.Err() != nil {
			queryErr = rows.Err()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := st.db.NamedQuery(fmt.Sprintf(sqlQueryKeysCounter, where), filter)
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
