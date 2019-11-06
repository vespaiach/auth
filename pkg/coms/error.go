package coms

import (
	"fmt"
)

type ErrCode int

const (
	// ErrMissingCredentials is to indicate that credentials is missing
	ErrMissingCredentials ErrCode = iota

	// ErrInvalidCredentials is to indicate that credentials is not valid
	ErrInvalidCredentials

	// ErrPermissionDenied is to indicate that wrong permission provided
	ErrPermissionDenied

	// ErrTokenExpired is to indicate that provided token expired
	ErrTokenExpired
)

type AppErr struct {
	err  error
	Code ErrCode
}

func NewAppErr(err error, code ErrCode) *AppErr {
	return &AppErr{err, code}
}

// Error func returns error message
func (appErr *AppErr) Error() string {
	return fmt.Sprintf("%s; (%d)", appErr.err.Error(), appErr.Code)
}
