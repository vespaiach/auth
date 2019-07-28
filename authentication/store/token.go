package store

import (
	"time"
)

// IssueTokenAction const
const IssueTokenAction string = "issue_token"

// RefreshTokenAction const
const RefreshTokenAction string = "refresh_token"

// Token model
type Token struct {
	ID           string `gorm:"type:varchar(16);primary_key"`
	UserID       uint   `gorm:"type:varchar(255)"`
	AccessToken  string `gorm:"type:varchar(2047)"`
	RefreshToken string `gorm:"type:varchar(2047)"`
	Action       string `gorm:"type:varchar(20)"`
	CreatedAt    time.Time
}

func (sto *mysqlStore) SaveToken(id string, userID uint, accessToken string, refreshToken string, action string) error {

	token := Token{
		ID:           id,
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Action:       action,
	}

	sto.db.Create(&token)

	if sto.db.NewRecord(token) {
		return ErrCanNotCreateData
	}

	return nil
}
