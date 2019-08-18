package model

import (
	"time"
)

// UserAction model
type UserAction struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	User      *User     `json:"user,omitempty"`
	ActionID  int64     `json:"action_id" db:"action_id"`
	Action    *Action   `json:"action,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserActionRepo defines user-action repo
type UserActionRepo interface {
	// GetByID gets user-action by ID
	GetByID(id int64) (*UserAction, error)

	// Create a new user-action
	Create(userID int64, roleID int64) (int64, error)

	// Delete user-action
	Delete(id int64) error

	// Query a list of user-actions
	Query(page int, perPage int, filters map[string]interface{}) ([]*UserAction, int64, error)
}
