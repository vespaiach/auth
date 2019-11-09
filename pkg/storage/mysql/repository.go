package mysql

import "github.com/jmoiron/sqlx"

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
