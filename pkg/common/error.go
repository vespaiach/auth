package common

import "errors"

type ErrCode int

const (
	ErrGetData ErrCode = iota
	ErrExecData
	ErrDataNotFound
	ErrDataFailValidation
)

var (
	ErrDuplicatedKey      = errors.New("duplicated key")
	ErrKeyNameInvalid     = errors.New("key name is invalid")
	ErrKeyNotFound        = errors.New("key doesn't exist")
	ErrBunchNotFound      = errors.New("bunch doesn't exist")
	ErrWrongInputDatatype = errors.New("inputted data type is incorrect")
	ErrDuplicatedBunch    = errors.New("duplicated bunch")
	ErrBunchNameInvalid   = errors.New("bunch name is invalid")
)
