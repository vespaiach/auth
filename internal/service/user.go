package service

import (
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// UserService is the interface that provides user's methods.
type UserService interface {

	// VerifyLogin method will return user entity
	VerifyLogin(username string, password string) (*model.User, error)

	// IssueTokens method will issue a pair of token
	IssueTokens(user *model.User) (*Credential, error)

	// RegisterUser method to add new user
	RegisterUser(fullName string, username string, password string, email string) (*model.User, error)
}

// Credential is model for user's authentication
type Credential struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type userService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// NewUserService creates a struct that implement IUserService
func NewUserService(appRepo *model.AppRepo, appConfig *conf.AppConfig) UserService {
	return &userService{
		appConfig,
		appRepo,
	}
}

func (s *userService) VerifyLogin(username string, password string) (user *model.User, err error) {
	repo := s.appRepo.UserRepo

	user, err = repo.GetByUsername(username)

	if user != nil && !isPasswordMatched(user.Hashed, password) {
		err = comtype.ErrCredentialNotMatch
		user = nil
	}

	return
}

func (s *userService) IssueTokens(u *model.User) (*Credential, error) {
	// id := uuid.New().String()
	// rsaConfig := s.appConfig.RsaKeyConfig
	// tokenConfig := s.appConfig.TokenConfig
	// repo := s.appRepo.UserRepo
	// tkrepo := s.appRepo.TokenRepo

	// accessToken, err := issueAccessToken(u, id, tokenConfig.AccessTokenDuration)
	// if err != nil {
	// 	return nil, err
	// }

	// var refreshToken string

	// if tokenConfig.UseRefreshToken {
	// 	refreshToken, err = issueRefreshToken(u, id, tokenConfig.RefreshTokenDuration)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// go tkrepo.Save(id, u.ID, accessToken, refreshToken)

	// return &Credential{
	// 	accessToken,
	// 	refreshToken,
	// }, nil
	return nil, nil
}

func (s *userService) RegisterUser(fullName string, username string, password string, email string) (user *model.User, err error) {
	comConfig := s.appConfig.CommonConfig
	repo := s.appRepo.UserRepo

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), comConfig.BcryptCost)
	if err != nil {
		return
	}

	user, err = repo.Create(fullName, username, string(hashedPassword), email)
	if err != nil {
		user = nil
	}

	return
}

func isPasswordMatched(hashed string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
