package modifying

import (
	"database/sql"
	"errors"
	"regexp"
)

// User model
type User struct {
	ID           int64
	Username     string
	Email        string
	Active       sql.NullBool
}

var ErrUserIDMissing = errors.New("user id is missing")
var ErrUsernameTooLong = errors.New("username exceeds 32 characters")
var ErrUsernameInvalid = errors.New("username contains special characters or white space characters")
var ErrEmailTooLong = errors.New("email address exceeds 64 characters")
var ErrEmailInvalid = errors.New("email address is invalid")
var ErrDuplicatedUsername = errors.New("username is duplicated")
var ErrDuplicatedEmail = errors.New("email address is duplicated")
var ErrPasswordHashedRequired = errors.New("password hash is missing")

func (u *User) Validate() error {
	if u.ID == 0 {
		return ErrUserIDMissing
	}

	if len(u.Username) > 32 {
		return ErrUsernameTooLong
	}

	if len(u.Username) > 0 {
		matched, err := regexp.Match(`^[a-z0-9_]+$`, []byte(u.Username))
		if !matched || err != nil {
			return ErrUsernameInvalid
		}
	}

	if len(u.Email) > 32 {
		return ErrEmailTooLong
	}

	if len(u.Email) > 0 {
		matched, err := regexp.Match(`^[a-z0-9_@\\-\\.]{1,127}$`, []byte(u.Email))
		if !matched || err != nil {
			return ErrEmailInvalid
		}
	}

	return nil
}
