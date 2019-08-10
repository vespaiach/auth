package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// Role model
type Role struct {
	ID        int64    `json:"id"`
	RoleName  string    `json:"role_name"`
	RoleDesc  string    `json:"role_desc"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Actions   []*Action `json:"actions"`
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
	Query(page int64, perPage int64, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*Role, int64, error)
}
