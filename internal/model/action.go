package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// Action model
type Action struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	ActionName string    `gorm:"type:varchar(63);unique_index" json:"action_name"`
	ActionDesc string    `gorm:"type:varchar(255)" json:"action_desc"`
	Active     bool      `gorm:"type:tinyint(1);index;default:1" json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ActionRepo defines action repo
type ActionRepo interface {
	// GetByID gets action by action ID
	GetByID(id uint) (*Action, error)

	// GetByName gets action by action's name
	GetByName(name string) (*Action, error)

	// Create a new action
	Create(name string, desc string) (*Action, error)

	// Update action
	Update(id uint, fields map[string]interface{}) error

	// Query a list of actions
	Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*Action, int64, error)
}
