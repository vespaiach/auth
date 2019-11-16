package usrmgr

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Hash      string
	Active    sql.NullBool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Bunch struct {
	ID        int64
	Name      string
	Desc      string
	Active    sql.NullBool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Key struct {
	ID        int64
	Key       string
	Desc      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
