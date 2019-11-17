package common

import "errors"

var (
	ErrDuplicatedKey      = errors.New("duplicated key")
	ErrKeyNameInvalid     = errors.New("key name is invalid")
	ErrKeyNotFound        = errors.New("key doesn't exist")
	ErrBunchNotFound      = errors.New("bunch doesn't exist")
	ErrWrongInputDatatype = errors.New("inputted data type is incorrect")
	ErrDuplicatedBunch    = errors.New("duplicated bunch")
	ErrBunchNameInvalid   = errors.New("bunch name is invalid")
	ErrUsernameInvalid    = errors.New("username is invalid")
	ErrWrongCredentials   = errors.New("wrong username or password")
	ErrDuplicatedUsername = errors.New("duplicated username")
	ErrEmailInvalid       = errors.New("email is invalid")
	ErrDuplicatedEmail    = errors.New("duplicated email")
	ErrMissingHash        = errors.New("hash is missing")
	ErrUserNotFound       = errors.New("user doesn't exist")
	ErrPasswordMissing    = errors.New("password is missing")
	ErrMissingJWTToken    = errors.New("jwt token is missing")
	ErrWrongJWTToken      = errors.New("jwt token is not correct")
	ErrNotAllowed         = errors.New("not allowed to access")
)
