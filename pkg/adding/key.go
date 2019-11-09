package adding

import "errors"

// ServiceKey model
type ServiceKey struct {
	Key string `json:"key" db:"key"`
	Desc string `json:"desc" db:"desc"`
}

var ErrKeyNameInvalid = errors.New("key name must be from 1 to 32 characters")
var ErrDuplicatedKey = errors.New("key name is duplicated")
var ErrKeyDescTooLong = errors.New("key description must be less than 64 characters")


func (sk *ServiceKey) Validate() error {
	if len(sk.Key) == 0 || len(sk.Key) > 32 {
		return ErrKeyNameInvalid
	}

	if len(sk.Desc) > 64 {
		return ErrKeyDescTooLong
	}

	return nil
}
