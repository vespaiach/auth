package models

import (
	"time"
)

// Role model
type Role struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	RoleName  string    `gorm:"type:varchar(63);unique_index" json:"role_name"`
	RoleDesc  string    `gorm:"type:varchar(255)" json:"role_desc"`
	Active    bool      `gorm:"type:tinyint(1);index;default:1" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
