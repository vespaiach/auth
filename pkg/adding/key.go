package adding

import (
	"errors"
	"regexp"
)

// ServiceKey model
type ServiceKey struct {
	Key  string
	Desc string
}

var ErrKeyInvalid = errors.New("key name must be from 1 to 32 characters")
var ErrKeyDescTooLong = errors.New("key description must be less than 64 characters")

func (sk *ServiceKey) Validate() error {
	if matched, err := regexp.Match(`^[a-z0-9_]{1,32}$`, []byte(sk.Key)); !matched || err != nil {
		return ErrKeyInvalid
	}

	if len(sk.Desc) > 64 {
		return ErrKeyDescTooLong
	}

	return nil
}
