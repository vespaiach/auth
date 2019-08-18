package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// TokenHistory model
type TokenHistory struct {
	UID         string    `json:"uid" db:"uid"`
	UserID      int64     `json:"user_id" db:"user_id"`
	AccessToken string    `json:"access_token" db:"access_token"`
	RefeshToken string    `json:"refresh_token" db:"refresh_token"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// TokenHistoryRepo defines token-history repo
type TokenHistoryRepo interface {
	// GetByID gets token-history by role ID
	GetByID(id int64) (*TokenHistory, error)

	// Create a new history
	Create(uid string, userID int64, accessToken string, refreshToken string, createdAt time.Time) error

	// Query a list of histories
	Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*TokenHistory, int64, error)
}
