package comtype

import "errors"

// Define all database errors
var (
	ErrDataNotFound    = errors.New("not found data")
	ErrCreadDataFailed = errors.New("can't create data")
)

// Define all datatype errors
var (
	ErrDataTypeMismatch = errors.New("mismatch data type")
	ErrNotAllowField    = errors.New("not allow field")
)

// Define all app common errors
var (
	ErrAppConfigMissingOrWrongSet = errors.New("app configurations are missing or wrong")
	ErrCredentialNotMatch         = errors.New("username or password is not valid")
)
