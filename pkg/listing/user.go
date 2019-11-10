package listing

import "time"

// User model
type User struct {
	ID        int64
	Username  string
	Email     string
	Hash      string
	Active    bool
	UpdatedAt time.Time
	CreatedAt time.Time
}