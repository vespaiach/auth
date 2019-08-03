package models

import (
	"time"
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
