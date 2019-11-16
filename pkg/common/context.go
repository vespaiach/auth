package common

// Context key
type ContextKey int

const (
	KeyManagementService ContextKey = iota
	AddingServiceContextKey
	ModifyServiceContextKey
	AppConfigContextKey
)
