package storage

import (
	"time"
)

//Key model
type Key struct {
	ID        int64
	Name      string
	Desc      string
	UpdatedAt time.Time
}

//QueryKey model
type QueryKey struct {
	Limit  int64
	Offset string
	Name   string
	Desc   string
	From   time.Time
	To     time.Time
}

//QueryKey model
type SortKey struct {
	Name      Direction
	Desc      Direction
	UpdatedAt Direction
}

//KeyStorer defines fundamental functions to interact with storage repository
type KeyStorer interface {
	Insert(k Key) (int64, error)
	Update(k Key) error
	Delete(id int64) error

	Get(id int64) (*Key, error)
	Query(queries QueryKey, sorts SortKey) ([]*Key, error)
}
