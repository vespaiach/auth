package coms

// ResponseCode is final result code
type ResponseCode string

// Define response codes
const (
	// ResponseSuccess is to indicate that request is successful
	ResponseSuccess ResponseCode = "success"

	// ResponseValidationFail is to indicate that request's validation fail
	ResponseValidationFail ResponseCode = "fail_validation"

	// ResponseForbiddenParams is to indicate that request has some forbidden params
	ResponseForbiddenParams ResponseCode = "forbid_param"

	// ResponseMissingParams is to indicate that request is missing some params
	ResponseMissingParams ResponseCode = "missing_param"

	// ResponseInvalidParams is to indicate that some params are invalid
	ResponseInvalidParams ResponseCode = "invalid_param"

	// 	ResponseRequireCredentials is to indicate that credentials are required
	ResponseRequireCredentials ResponseCode = "require_credentials"

	// ResponseWrongPermission is to indicate that provided permission is wrong
	ResponseWrongPermission ResponseCode = "wrong_permission"

	// ResponseServerError is to indicate that server encounter error
	ResponseServerError ResponseCode = "server_error"

	// ResponseDataNotFound is to indicate that data was not found
	ResponseDataNotFound ResponseCode = "not_found_data"
)
