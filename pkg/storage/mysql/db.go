package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"sync"
)

var initDatabase = `
CREATE TABLE IF NOT EXISTS "keys" (
  	"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	"key" VARCHAR(32) NOT NULL DEFAULT '',
	"desc" VARCHAR(64) NOT NULL DEFAULT '',
  	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  	"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "keys_key_uniq" ("key" ASC))
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "bunches" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "name" VARCHAR(32) NOT NULL DEFAULT '',
  "desc" VARCHAR(64) NOT NULL DEFAULT '',
  "active" TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "bunch_name_uniq" ("name" ASC),
  INDEX "bunch_active_idx" ("active" ASC))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS "users" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "username" VARCHAR(32) NOT NULL,
  "email" VARCHAR(64) NOT NULL,
  "hash" VARCHAR(255) NOT NULL,
  "active" TINYINT(1) NOT NULL DEFAULT 1,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "users_username_uniq" ("username" ASC),
  UNIQUE INDEX "users_email_uniq" ("email" ASC),
  INDEX "users_active_idx" ("active" ASC))
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "token_histories" (
  "uid" VARCHAR(36) NOT NULL,
  "user_id" BIGINT(20) UNSIGNED NOT NULL,
  "access_token" VARCHAR(1024) NOT NULL,
  "refresh_token" VARCHAR(1024) NOT NULL DEFAULT '',
  "remote_addr" VARCHAR(512) NOT NULL DEFAULT '',
  "x_forwarded_for" VARCHAR(512) NOT NULL DEFAULT '',
  "x_real_ip" VARCHAR(512) NOT NULL DEFAULT '',
  "user_agent" VARCHAR(512) NOT NULL DEFAULT '',
  "created_at" TIMESTAMP NOT NULL,
  "expired_at" TIMESTAMP NOT NULL,
  PRIMARY KEY ("uid"),
  UNIQUE INDEX "uid_uniq" ("uid" ASC))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "bunch_keys" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "bunch_id" BIGINT(20) UNSIGNED NOT NULL,
  "key_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "bunch_key_key_id_idx" ("key_id" ASC),
  INDEX "bunch_key_bunch_id_idx" ("bunch_id" ASC),
  UNIQUE INDEX "bunch_key_uniq" ("bunch_id" ASC, "key_id" ASC),
  CONSTRAINT "key_id_on_bunch_key"
    FOREIGN KEY ("key_id")
    REFERENCES "keys" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "role_id_on_bunch_key"
    FOREIGN KEY ("bunch_id")
    REFERENCES "bunches" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "user_bunches" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "user_id" BIGINT(20) UNSIGNED NOT NULL,
  "bunch_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "user_bunch_user_id_idx" ("user_id" ASC),
  INDEX "user_bunch_bunch_id_idx" ("bunch_id" ASC),
  UNIQUE INDEX "user_bunch_uniq" ("user_id" ASC, "bunch_id" ASC),
  CONSTRAINT "user_id_on_user_bunch"
    FOREIGN KEY ("user_id")
    REFERENCES "users" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "bunch_id_on_user_bunch"
    FOREIGN KEY ("bunch_id")
    REFERENCES "bunches" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB;
`

var dropDatabase = `
DROP TABLE IF EXISTS "user_bunches";
DROP TABLE IF EXISTS "bunch_keys";
DROP TABLE IF EXISTS "keys";
DROP TABLE IF EXISTS "bunches";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "token_histories";
`

type uniqInt struct {
	order int
	mux   sync.Mutex
}

var inc *uniqInt = &uniqInt{order: 1}

func (unq *uniqInt) New() int {
	unq.mux.Lock()
	defer unq.mux.Unlock()
	unq.order++
	return unq.order
}

// createUniqueString is to create unique string for testing
func (m *Migrator) createUniqueString(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, strconv.Itoa(inc.New()))
}

// Script migration script
type Script struct {
	Name string
	Text string
}

// Migrator struct
type Migrator struct {
	db   *sqlx.DB
	init []*Script
	drop []*Script
}

// NewMigrator return struct instance
func NewMigrator(db *sqlx.DB) *Migrator {

	var initScripts = []*Script{
		&Script{Name: "init_database", Text: initDatabase},
	}

	var dropScripts = []*Script{
		&Script{Name: "drop_database", Text: dropDatabase},
	}

	return &Migrator{
		db,
		initScripts,
		dropScripts,
	}
}

// Init database
func (m *Migrator) Init() {
	tx := m.db.MustBegin()

	for _, s := range m.init {
		tx.MustExec(santizeSQL(s.Text))
	}

	tx.Commit()
}

// Drop database
func (m *Migrator) Drop() {
	tx := m.db.MustBegin()

	for i := len(m.drop) - 1; i >= 0; i-- {
		tx.MustExec(santizeSQL(m.drop[i].Text))
	}

	tx.Commit()
}

func santizeSQL(sql string) string {
	return strings.Replace(sql, `"`, "`", -1)
}

func (m *Migrator) createSeedingServiceKey(beforeCreate func(map[string]interface{})) int64 {
	fields := map[string]interface{}{
		"key":  fmt.Sprintf("key_%s", strconv.Itoa(inc.New())),
		"desc": fmt.Sprintf("desc_%s", strconv.Itoa(inc.New())),
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec("INSERT INTO `keys` (`key`, `desc`) VALUES (:key, :desc);", fields)
	id, _ := result.LastInsertId()

	return id
}

func (m *Migrator) createSeedingBunch(beforeCreate func(map[string]interface{})) int64 {
	fields := map[string]interface{}{
		"name":   fmt.Sprintf("name_%s", strconv.Itoa(inc.New())),
		"desc":   fmt.Sprintf("desc_%s", strconv.Itoa(inc.New())),
		"active": true,
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec("INSERT INTO `bunches` (`name`, `desc`, active) VALUES (:name, :desc, :active);",
		fields)
	id, _ := result.LastInsertId()

	return id
}

func (m *Migrator) createSeedingUser(beforeCreate func(map[string]interface{})) int64 {
	fields := map[string]interface{}{
		"username": fmt.Sprintf("username_%s", strconv.Itoa(inc.New())),
		"email":    fmt.Sprintf("email_%s", strconv.Itoa(inc.New())),
		"hash":     fmt.Sprintf("hash_%s", strconv.Itoa(inc.New())),
		"active":   true,
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec("INSERT INTO users(username, email, hash, active) "+
		"VALUES(:username, :email, :hash, :active);", fields)
	id, _ := result.LastInsertId()

	return id
}

func (m *Migrator) getServiceKeyByID(id int64) (key string, desc string) {
	rows, err := m.db.Queryx("SELECT `key`, `desc` FROM `keys` WHERE id = ?", id)
	defer rows.Close()

	if err == nil && rows.Next() {
		rows.Scan(&key, &desc)
	}

	return
}

func (m *Migrator) getBunchByID(id int64) (name string, desc string, active bool) {
	rows, err := m.db.Queryx("SELECT `name`, `desc`, `active` FROM bunches WHERE id = ?", id)
	defer rows.Close()

	if err == nil && rows.Next() {
		rows.Scan(&name, &desc, &active)
	}

	return
}

func (m *Migrator) getUserByID(id int64) (username string, email string, hash string, active bool) {
	rows, err := m.db.Queryx("Select `username`, `email`, `hash`, `active` FROM `users` WHERE id = ?", id)
	defer rows.Close()

	if err == nil && rows.Next() {
		rows.Scan(&username, &email, &hash, &active)
	}

	return
}

func (m *Migrator) getKeyIDByBunchID(id int64) []int64 {
	rows, _ := m.db.Queryx("SELECT key_id FROM `bunch_keys` WHERE bunch_id = ?", id)
	defer rows.Close()

	results := make([]int64, 0)

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		results = append(results, id)
	}

	return results
}
