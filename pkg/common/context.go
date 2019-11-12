package common

// Context key
type ContextKey int

const (
	ListingServiceContextKey ContextKey = iota
	AddingServiceContextKey
	AppConfigContextKey
)
