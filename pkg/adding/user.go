package adding

import (
	"errors"
	"regexp"
)

// User model
type User struct {
	Username string
	Email    string
	Hash     string
}

var ErrUsernameRequired = errors.New("username is missing")
var ErrUsernameTooLong = errors.New("username exceeds 32 characters")
var ErrUsernameInvalid = errors.New("username contains special characters or white space characters")
var ErrEmailRequired = errors.New("email address is missing")
var ErrEmailTooLong = errors.New("email address exceeds 64 characters")
var ErrEmailInvalid = errors.New("email address is invalid")
var ErrHashRequired = errors.New("password hash is missing")

// Validate user data before adding
func (u *User) Validate() error {
	if len(u.Username) == 0 {
		return ErrUsernameRequired
	}

	if len(u.Username) > 32 {
		return ErrUsernameTooLong
	}

	if matched, err := regexp.Match(`^[a-zA-Z0-9_]{1,32}$`, []byte(u.Username)); !matched || err != nil {
		return ErrUsernameInvalid
	}

	if len(u.Email) == 0 {
		return ErrEmailRequired
	}

	if len(u.Email) > 32 {
		return ErrEmailTooLong
	}

	if matched, err := regexp.Match(`^[a-zA-Z0-9_@\\-\\.]{1,127}$`, []byte(u.Email)); !matched || err != nil {
		return ErrEmailInvalid
	}

	if len(u.Hash) == 0 {
		return ErrHashRequired
	}

	return nil
}
