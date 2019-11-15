package modifying

import (
	"errors"
)

// ServiceKey model
type ServiceKey struct {
	ID   int64
	Key  string
	Desc string
}

var ErrServiceKeyIDMissing = errors.New("key id is missing")
var ErrServiceKeyTooLong = errors.New("key name is too long")
var ErrServiceKeyDescTooLong = errors.New("key desc is too long")
var ErrDuplicatedKey = errors.New("key name is duplicated")

func (sk *ServiceKey) Validate() error {
	if sk.ID == 0 {
		return ErrServiceKeyIDMissing
	}

	if len(sk.Key) > 32 {
		return ErrServiceKeyTooLong
	}

	if len(sk.Desc) > 64 {
		return ErrServiceKeyDescTooLong
	}

	return nil
}
