package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// User model
type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	FullName  string    `gorm:"type:varchar(255)" json:"full_name"`
	Username  string    `gorm:"type:varchar(63);unique_index" json:"username"`
	Hashed    string    `gorm:"type:varchar(511)"`
	Email     string    `gorm:"type:varchar(255);unique_index" json:"email"`
	Active    int       `gorm:"type:int,index" json:"active"`
	Verified  bool      `gorm:"type:bool,index" json:"verified"`
	Actions   []*Action `gorm:"many2many:user_actions;"`
	Roles     []*Role   `gorm:"many2many:user_roles;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepo defines user repo
type UserRepo interface {
	// GetByID gets user by user ID
	GetByID(id uint) (user *User, err error)

	// GetByEmail gets user by user's email
	GetByUsername(username string) (*User, error)

	// GetByEmail gets user by user's email
	GetByEmail(email string) (*User, error)

	// Create a new user
	Create(fullName string, username string, hashed string, email string) (*User, error)

	// Update user
	Update(id uint, fields map[string]interface{}) error

	// Query a list of users
	Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*User, int64, error)
}
