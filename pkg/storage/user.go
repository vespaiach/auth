package storage

import (
	"database/sql"
	"time"
)

//User model
type User struct {
	ID        int64
	FullName  string
	Username  string
	Email     string
	Hash      string
	Salt      string
	Active    sql.NullBool
	UpdatedAt time.Time
}

//UserBunch model
type UserBunch struct {
	ID      int64
	UserID  int64
	BunchID int64
}

//UserStorer defines fundamental functions to interact with storage repository
type UserStorer interface {
	Insert(u User) (*User, error)
	Update(u User) (*User, error)
	Delete(id int64) error

	Get(id int64) (*User, error)
	Query(limit int64, offset int64, rules map[string]interface{}, sorts map[string]interface{}) ([]*User, error)
}

//UserBunchStorer defines fundamental functions to interact with storage repository
type UserBunchStorer interface {
	Insert(ub UserBunchStorer) (*UserBunchStorer, error)
	Update(ub UserBunchStorer) (*UserBunchStorer, error)
	Delete(id int64) error

	Get(id int64) (*UserBunchStorer, error)
	Query(limit int64, offset int64, rules map[string]interface{},
		sorts map[string]interface{}) ([]*UserBunchStorer, int64, error)
}
