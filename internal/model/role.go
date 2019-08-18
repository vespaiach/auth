package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// Role model
type Role struct {
	ID        int64     `json:"id" db:"id"`
	RoleName  string    `json:"role_name" db:"role_name"`
	RoleDesc  string    `json:"role_desc" db:"role_desc"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Actions   []*Action `json:"actions" db:"actions"`
}

// RoleRepo defines role repo
type RoleRepo interface {
	// GetByID gets role by role ID
	GetByID(id int64) (*Role, error)

	// GetByEmail gets role by role's email
	GetByName(rolename string) (*Role, error)

	// Create a new role
	Create(roleName string, roleDesc string) (int64, error)

	// Update role
	Update(id int64, fields map[string]interface{}) error

	// Query a list of roles
	Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*Role, int64, error)
}
