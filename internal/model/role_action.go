package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// RoleAction model
type RoleAction struct {
	ID        int64     `json:"id" db:"id"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	Role      *Role     `json:"role,omitempty"`
	ActionID  int64     `json:"action_id" db:"action_id"`
	Action    *Action   `json:"action,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RoleActionRepo defines role-action repo
type RoleActionRepo interface {
	// GetByID gets role-action by ID
	GetByID(id int64) (*RoleAction, *comtype.CommonError)

	// Create a new role-action
	Create(roleID int64, actionID int64) (int64, *comtype.CommonError)

	// Delete role-action
	Delete(id int64) *comtype.CommonError

	// Query a list of role-actions
	Query(take int, filters map[string]interface{}) ([]*RoleAction, *comtype.CommonError)
}
