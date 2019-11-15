package keymgr

import (
	"errors"
	"time"
)

type Key struct {
	ID        int64
	Key       string
	Desc      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrDuplicatedKey = errors.New("duplicated key")
var ErrKeyNameInvalid = errors.New("key name is invalid")
var ErrKeyNotFound = errors.New("key doesn't exist")
var ErrBunchNotFound = errors.New("bunch doesn't exist")
