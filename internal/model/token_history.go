package model

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// TokenHistory model
type TokenHistory struct {
	UID           string    `json:"uid" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	AccessToken   string    `json:"access_token" db:"access_token"`
	RefreshToken  string    `json:"refresh_token" db:"refresh_token"`
	RemoteAddr    string    `json:"remote_addr" db:"remote_addr"`
	XForwardedFor string    `json:"x_forwarded_for" db:"x_forwarded_for"`
	XRealIP       string    `json:"x_real_ip" db:"x_real_ip"`
	UserAgent     string    `json:"user_agent" db:"user_agent"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	ExpiredAt     time.Time `json:"expired_at" db:"expired_at"`
}

// TokenHistoryRepo defines token repo
type TokenHistoryRepo interface {
	// Save a TokenHistory
	Save(uid string, userID int64, accessToken string, refreshToken string, remoteAddr string,
		xForwardedFor string, xRealIP string, userAgent string, createdAt time.Time, expiredAt time.Time) *comtype.CommonError

	// GetByUserID find all isssued tokens by user's ID
	GetByUserID(userID int64) ([]*TokenHistory, *comtype.CommonError)
}
