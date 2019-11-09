package adding

import "errors"

// Bunch model
type Bunch struct {
	Name string `json:"name" db:"name"`
	Desc string `json:"desc" db:"desc"`
}

var ErrBunchNameInvalid = errors.New("bunch name must be from 1 to 32 characters")
var ErrDuplicatedBunch = errors.New("bunch name is duplicated")
var ErrBunchDescTooLong = errors.New("bunch description must be less than 64 characters")

func (b *Bunch) Validate() error {
	if len(b.Name) == 0 || len(b.Name) > 32 {
		return ErrBunchNameInvalid
	}

	if len(b.Desc) > 64 {
		return ErrBunchDescTooLong
	}

	return nil
}
