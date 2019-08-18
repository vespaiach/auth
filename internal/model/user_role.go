package model

import (
	"time"
)

// UserRole model
type UserRole struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	User      *User     `json:"user,omitempty"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	Role      *Role     `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserRoleRepo defines user-role repo
type UserRoleRepo interface {
	// GetByID gets user-role by ID
	GetByID(id int64) (*UserRole, error)

	// Create a new user-role
	Create(userID int64, roleID int64) (int64, error)

	// Delete user-role
	Delete(id int64) error

	// Query a list of user-roles
	Query(page int, perPage int, filters map[string]interface{}) ([]*UserRole, int64, error)
}
