package adding

import (
	"errors"
	"regexp"
)

// Bunch model
type Bunch struct {
	Name string
	Desc string
}

var ErrBunchInvalid = errors.New("bunch name must be from 1 to 32 characters")
var ErrBunchDescTooLong = errors.New("bunch description must be less than 64 characters")

func (b *Bunch) Validate() error {
	if matched, err := regexp.Match(`^[a-z0-9_]{1,32}$`, []byte(b.Name)); !matched || err != nil {
		return ErrBunchInvalid
	}

	if len(b.Desc) > 64 {
		return ErrBunchDescTooLong
	}

	return nil
}
