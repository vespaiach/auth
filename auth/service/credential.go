package service

import (
	"github.com/google/uuid"
	"github.com/vespaiach/authentication/store"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) VerifyLogin(username string, password string) (user *store.User, err error) {
	user, err = s.store.GetUserByUsername(username)

	if user != nil && !isPasswordMatched(user.Hashed, password) {
		err = ErrCredentialNotMatch
		user = nil
	}

	return
}

func (s *service) IssueTokens(u *store.User) (*Credential, error) {
	id := uuid.New().String()

	accessToken, err := issueAccessToken(u, id, s.config)
	if err != nil {
		return nil, err
	}

	var refreshToken string

	if s.config.UseRefreshToken {
		refreshToken, err = issueRefreshToken(u, id, s.config)
		if err != nil {
			return nil, err
		}
	}

	go s.store.SaveToken(id, u.ID, accessToken, refreshToken, store.IssueTokenAction)

	return &Credential{
		accessToken,
		refreshToken,
	}, nil
}

func (s *service) RegisterUser(name string, username string, password string, email string) (user *store.User, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BcryptCost)
	if err != nil {
		return
	}

	user, err = s.store.CreateUser(name, username, string(hashedPassword), email)
	if err != nil {
		user = nil
	}

	return
}

func isPasswordMatched(hashed string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
