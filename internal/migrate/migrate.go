package migrate

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// Script migration script
type Script struct {
	Name string
	Text string
}

// Migrator struct
type Migrator struct {
	db   *sqlx.DB
	up   []*Script
	down []*Script
}

// NewMigrator return struct instance
func NewMigrator(db *sqlx.DB) *Migrator {

	var upScripts = []*Script{
		&Script{Name: "init_database", Text: initDatabase},
	}

	var downScripts = []*Script{
		&Script{Name: "drop_init_database", Text: dropInitDatabase},
	}

	return &Migrator{
		db:   db,
		up:   upScripts,
		down: downScripts,
	}
}

// Up to run migration scripts
func (m *Migrator) Up() {
	tx := m.db.MustBegin()

	for _, s := range m.up {
		tx.MustExec(santizeSQL(s.Text))
	}

	tx.Commit()
}

// Down to run migration scripts
func (m *Migrator) Down() {
	tx := m.db.MustBegin()

	for i := len(m.down) - 1; i >= 0; i-- {
		tx.MustExec(santizeSQL(m.down[i].Text))
	}

	tx.Commit()
}

// SeedTestData to migrate test data
func (m *Migrator) SeedTestData() {
	tx := m.db.MustBegin()
	tx.MustExec(santizeSQL(initTestData))
	tx.Commit()
}

func santizeSQL(sql string) string {
	return strings.Replace(sql, `"`, "`", -1)
}
