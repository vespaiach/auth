package common

// Context key
type ContextKey int

const (
	KeyManagementService ContextKey = iota
	BunchManagementService
	UserManagementService
	AppConfigContextKey
)
