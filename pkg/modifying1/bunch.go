package modifying

import (
	"database/sql"
	"errors"
)

// Bunch model
type Bunch struct {
	ID     int64
	Name   string
	Desc   string
	Active sql.NullBool
}

var ErrBunchIDMissing = errors.New("bunch id is missing")
var ErrDuplicatedBunch = errors.New("bunch name is duplicated")
var ErrBunchNameTooLong = errors.New("bunch name must be less than 32 characters")
var ErrBunchDescTooLong = errors.New("bunch description must be less than 64 characters")

func (b *Bunch) Validate() error {
	if b.ID == 0 {
		return ErrBunchIDMissing
	}

	if len(b.Name) > 32 {
		return ErrBunchNameTooLong
	}

	if len(b.Desc) > 64 {
		return ErrBunchDescTooLong
	}

	return nil
}
