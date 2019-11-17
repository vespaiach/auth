package common

import "time"

// TimeLayout common time format
const TimeLayout = time.RFC3339

// NullableBool introduce nullable type for boolean
type NullableBool struct {
	Value bool
	IsSet bool
}

// Default total record will be returned in a request to sql server
const Take = 100

// SortingDirection sort by direction
type SortingDirection int

const (
	Ascending SortingDirection = iota
	Descending
)
