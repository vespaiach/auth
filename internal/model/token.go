package model

import (
	"time"
)

// Token model
type Token struct {
	ID           string    `gorm:"primary_key" json:"id"`
	UserID       uint      `gorm:"type:int" json:"user_id"`
	AccessToken  string    `gorm:"type:varchar(255)" json:"full_name"`
	RefreshToken string    `gorm:"type:varchar(63);unique_index" json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TokenRepo defines token repo
type TokenRepo interface {
	// Save a token
	Save(id string, userID uint, accessToken string, refreshToken string) error

	// GetByID find token by ID
	GetByID(id string) (*Token, error)
}
