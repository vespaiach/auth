package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// User model
type User struct {
	ID        int64     `json:"id" db:"id"`
	FullName  string    `json:"full_name" db:"full_name"`
	Username  string    `json:"username" db:"username"`
	Hashed    string    `json:"-"`
	Email     string    `json:"email" db:"email"`
	Active    bool      `json:"active" db:"active"`
	Verified  bool      `json:"verified" db:"verified"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepo defines user repo
type UserRepo interface {
	// GetByID gets user by user ID
	GetByID(id int64) (*User, error)

	// GetByUsername gets user by user's username
	GetByUsername(username string) (*User, error)

	// GetByEmail gets user by user's email
	GetByEmail(email string) (*User, error)

	// Create a new user
	Create(fullName string, username string, hashed string, email string) (int64, error)

	// Update user
	Update(id int64, fields map[string]interface{}) error

	// Query a list of users
	Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*User, int64, error)
}
