package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/vespaiach/authentication/config"
	"github.com/vespaiach/authentication/store"
)

// ErrCredentialNotMatch username or password is wrong
var ErrCredentialNotMatch = errors.New("username or password is not valid")

// Service is the interface that provides auth methods.
type Service interface {

	// VerifyLogin method will return user entity
	VerifyLogin(username string, password string) (*store.User, error)

	// IssueTokens method will issue a pair of token
	IssueTokens(u *store.User) (*Credential, error)

	// RegisterUser method to add new user
	RegisterUser(name string, username string, password string, email string) (user *store.User, err error)
}

// Credential is model for user's authentication
type Credential struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type service struct {
	config *config.ServiceConfig
	store  store.Store
}

// NewService creates a auth service
func NewService(db *gorm.DB, config *config.ServiceConfig) Service {
	return &service{
		config: config,
		store:  store.NewStore(db, config),
	}
}
