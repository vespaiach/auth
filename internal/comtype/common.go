package comtype

import "time"

// ActivateStatus will have value: 0, 1, 2
type ActivateStatus int

// Defined activate status
const (
	Active   ActivateStatus = 1
	Unactive ActivateStatus = 0
	Unset    ActivateStatus = 2
)

// SortDirection vall have value 1,2
type SortDirection int

// Defined activate status
const (
	Ascending SortDirection = iota
	Decending
)

// Key is a general Key type
type Key string

// DateTimeLayout Layout
var DateTimeLayout = time.RFC3339

// ResponseCode is final result code
type ResponseCode string

// Defined result codes
const (
	SuccessCode               ResponseCode = "success"
	ValidationFailCode        ResponseCode = "fail_validation"
	ParamForbiddenCode        ResponseCode = "forbid_param"
	ParamMissingCode          ResponseCode = "missing_param"
	ParamInvalidCode          ResponseCode = "invalid_param"
	MissingAuthenticationCode ResponseCode = "missing_authentication"
	MissingPermissionCode     ResponseCode = "missing_permission"
	ServerErrorCode           ResponseCode = "server_error"
	NotFoundDataCode          ResponseCode = "not_found_data"
)

// TimeLayout common time format
var TimeLayout = time.RFC3339

type commonKey int

// Common keys
const (
	CommonKeyRequestContext commonKey = iota
	CommonKeyAppConfiguration
)
