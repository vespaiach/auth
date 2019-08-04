package comtype

// ActivateStatus will have value: 1, 2, 3
type ActivateStatus int

// Defined activate status
const (
	Active ActivateStatus = iota
	Unactive
	Unset
)

// SortDirection vall have value 1,2
type SortDirection int

// Defined activate status
const (
	Ascending SortDirection = iota
	Decending
)
