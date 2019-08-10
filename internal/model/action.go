package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// Action model
type Action struct {
	ID         int64     `json:"id" db:"id"`
	ActionName string    `json:"action_name" db:"action_name"`
	ActionDesc string    `json:"action_desc" db:"action_desc"`
	Active     bool      `json:"active" db:"active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ActionRepo defines action repo
type ActionRepo interface {
	// GetByID gets action by action ID
	GetByID(id int64) (*Action, error)

	// GetByName gets action by action's name
	GetByName(name string) (*Action, error)

	// Create a new action
	Create(name string, desc string) (int64, error)

	// Update action
	Update(id int64, fields map[string]interface{}) error

	// Query a list of actions
	Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*Action, int64, error)
}
