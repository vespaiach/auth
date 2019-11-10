package mysql

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
)

// Storage stores data in mysql server
type Storage struct {
	DbClient *sqlx.DB
}

// NewStorage create new instance of Storage
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db,
	}
}

// Init database connection
func InitDb(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
