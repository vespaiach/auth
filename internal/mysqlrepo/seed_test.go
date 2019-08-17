package mysqlrepo

import (
	"fmt"
	"strings"
	"sync"
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

func (t *appTesting) createActionTable() error {
	return t.runchema(`
		CREATE TABLE IF NOT EXISTS "actions" (
			"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			"action_name" VARCHAR(63) NOT NULL,
			"action_desc" VARCHAR(255) NOT NULL DEFAULT '',
			"active" TINYINT(1) NOT NULL DEFAULT 1,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY ("id"),
		UNIQUE INDEX "action_name_uniq" ("action_name" ASC),
		INDEX "action_active_idx" ("active" ASC))
		ENGINE = InnoDB
		AUTO_INCREMENT = 1
		DEFAULT CHARACTER SET = utf8;
	`)
}

func (t *appTesting) createRoleTable() error {
	return t.runchema(`
		CREATE TABLE IF NOT EXISTS "roles" (
			"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			"role_name" VARCHAR(63) NOT NULL DEFAULT '',
			"role_desc" VARCHAR(254) NOT NULL DEFAULT '',
			"active" TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY ("id"),
			UNIQUE INDEX "role_name_uniq" ("role_name" ASC),
			INDEX "role_active_idx" ("active" ASC))
		ENGINE = InnoDB;
	`)
}

func (t *appTesting) createUserTable() error {
	return t.runchema(`
		CREATE TABLE IF NOT EXISTS "users" (
			"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			"full_name" VARCHAR(254) NOT NULL DEFAULT '',
			"username" VARCHAR(63) NOT NULL DEFAULT '',
			"hashed" VARCHAR(255) NOT NULL DEFAULT '',
			"email" VARCHAR(127) NOT NULL DEFAULT '',
			"verified" TINYINT(1) NOT NULL DEFAULT 0,
			"active" TINYINT(1) NOT NULL DEFAULT 1,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY ("id"),
			UNIQUE INDEX "user_name_uniq" ("username" ASC),
			UNIQUE INDEX "user_email_uniq" ("email" ASC),
			INDEX "user_active_idx" ("active" ASC))
		ENGINE = InnoDB
		AUTO_INCREMENT = 1
		DEFAULT CHARACTER SET = utf8;
	`)
}

func (t *appTesting) dropActionTable() error {
	return t.runchema(`DROP TABLE actions;`)
}

func (t *appTesting) dropRoleTable() error {
	return t.runchema(`DROP TABLE roles;`)
}

func (t *appTesting) dropUserTable() error {
	return t.runchema(`DROP TABLE users;`)
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
