package mysqlrepo

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vespaiach/auth/internal/model"
)

var (
	mu                sync.RWMutex
	counter           = 0
	defaultFixtureRow = 20
)

func (t *appTesting) santizeSQL(sql string) string {
	return strings.Replace(sql, `"`, "`", -1)
}

func (t *appTesting) runchema(sql string) error {
	sql = t.santizeSQL(sql)

	stmts := strings.Split(sql, ";\n")
	if len(strings.Trim(stmts[len(stmts)-1], " \n\t\r")) == 0 {
		stmts = stmts[:len(stmts)-1]
	}
	for _, s := range stmts {
		_, err := t.db.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *appTesting) generateUniqueString(prefix string) string {
	var str string
	mu.Lock()
	counter = counter + 1
	str = fmt.Sprintf("%s_uniq_%d", prefix, counter)
	mu.Unlock()

	return str
}

func (t *appTesting) loadRoleFixtures(rolePrefix string) ([]int64, error) {
	ids := make([]int64, 0, defaultFixtureRow)

	tx, _ := t.db.Begin()
	stmt, err := tx.Prepare(sqlCreateRole)
	if err != nil {
		return nil, err
	}

	for i := 0; i < defaultFixtureRow; i++ {
		name := t.generateUniqueString(rolePrefix)
		res, err := stmt.Exec(name, name)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		ids = append(ids, lastID)
	}
	tx.Commit()

	return ids, nil
}

func (t *appTesting) loadActionFixtures(actionPrefix string) ([]int64, error) {
	ids := make([]int64, 0, defaultFixtureRow)

	tx, _ := t.db.Begin()
	stmt, err := tx.Prepare(sqlCreateAction)
	if err != nil {
		return nil, err
	}

	for i := 0; i < defaultFixtureRow; i++ {
		name := t.generateUniqueString(actionPrefix)
		res, err := stmt.Exec(name, name)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		ids = append(ids, lastID)
	}
	tx.Commit()

	return ids, nil
}

func (t *appTesting) loadUserFixtures(usernamePrefix string) ([]int64, error) {
	ids := make([]int64, 0, defaultFixtureRow)

	tx, _ := t.db.Begin()
	stmt, err := tx.Prepare(sqlCreateUser)
	if err != nil {
		return nil, err
	}

	for i := 0; i < defaultFixtureRow; i++ {
		name := t.generateUniqueString(usernamePrefix)
		res, err := stmt.Exec(name, name, name, name)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		ids = append(ids, lastID)
	}
	tx.Commit()

	return ids, nil
}

func (t *appTesting) createActionWithName(name string) (int64, error) {
	stmt, err := t.db.Prepare(sqlCreateAction)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(name, name)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *appTesting) createRoleWithName(name string) (int64, error) {
	stmt, err := t.db.Prepare(sqlCreateRole)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(name, name)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *appTesting) createUserWithName(name string) (int64, error) {
	stmt, err := t.db.Prepare(sqlCreateUser)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(name, name, name, name)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *appTesting) createTokenHistory(beforeSave func(*model.TokenHistory)) error {
	history := &model.TokenHistory{
		UID:         uuid.New().String(),
		UserID:      1,
		AccessToken: "Access Token",
		CreatedAt:   time.Now(),
	}

	if beforeSave != nil {
		beforeSave(history)
	}

	res, err := t.db.NamedExec(`
	INSERT INTO token_histories (uid, user_id, access_token, refresh_token, created_at) 
	VALUES (:uid, :user_id, :access_token, :refresh_token, :created_at)`, history)
	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		return errors.New("couldn't create token-history")
	}

	return nil
}
