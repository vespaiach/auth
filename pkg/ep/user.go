package ep

import (
	"context"
	"errors"
	"regexp"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/listing"
	"github.com/vespaiach/auth/pkg/modifying"
	"golang.org/x/crypto/bcrypt"
)

// AddingUser model for adding a user
type AddingUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdatingUser model for updatina user
type UpdatingUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token model
type Token struct {
	AccessToken string `json:"access_token"`
}

type tokenClaims struct {
	Bunches []string
	Keys    []string
	jwtgo.StandardClaims
}

func (u *UpdatingUser) Validate() error {
	payload := make(map[string]string)

	if matched, err := regexp.Match(`^[a-zA-Z0-9_]{1,32}$`, []byte(u.Username)); !matched || err != nil {
		payload["username"] = "username is not valid"
	}

	if matched, err := regexp.Match(`^[a-zA-Z0-9_@\\-\\.]{1,127}$`, []byte(u.Email)); !matched || err != nil {
		payload["email"] = "email is not valid"
	}

	if len(u.Password) == 0 {
		payload["password"] = "password is required"
	}

	if len(payload) > 0 {
		err := common.NewAppErr(errors.New("data is not valid"), common.ErrDataFailValidation)
		err.Payload = payload
		return err
	}

	return nil
}

func (u *AddingUser) Validate() error {
	payload := make(map[string]string)

	if matched, err := regexp.Match(`^[a-zA-Z0-9_]{1,32}$`, []byte(u.Username)); !matched || err != nil {
		payload["username"] = "username is not valid"
	}

	if matched, err := regexp.Match(`^[a-zA-Z0-9_@\\-\\.]{1,127}$`, []byte(u.Email)); !matched || err != nil {
		payload["email"] = "email is not valid"
	}

	if len(u.Password) == 0 {
		payload["password"] = "password is required"
	}

	if len(payload) > 0 {
		err := common.NewAppErr(errors.New("data is not valid"), common.ErrDataFailValidation)
		err.Payload = payload
		return err
	}

	return nil
}

// VerifyingUser model for verifying a user
type VerifyingUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AddUser endpoint is to make an endpoint for adding a new user
func AddUser(ctx context.Context, request interface{}) (interface{}, error) {
	addingServ := ctx.Value(common.AddingServiceContextKey).(adding.Service)
	listingServ := ctx.Value(common.ListingServiceContextKey).(listing.Service)
	appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

	addingUser, ok := request.(*AddingUser)
	if !ok {
		return nil, errors.New("couldn't get user data")
	}

	if errValidation := addingUser.Validate(); errValidation != nil {
		return nil, errValidation
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(addingUser.Password), appConfig.BcryptCost)
	if err != nil {
		return nil, err
	}

	id, err := addingServ.AddUser(adding.User{
		Username: addingUser.Username,
		Email:    addingUser.Email,
		Hash:     string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	user, err := listingServ.GetUser(id)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &User{
			user.ID,
			user.Username,
			user.Email,
			user.Active,
			user.CreatedAt,
			user.UpdatedAt,
		}, nil
	}
}

// VerifyUser endpoint is to verify a existing user
func VerifyUser(ctx context.Context, request interface{}) (interface{}, error) {
	listingServ := ctx.Value(common.ListingServiceContextKey).(listing.Service)
	appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

	verifying, ok := request.(*VerifyingUser)
	if !ok {
		return nil, errors.New("couldn't get user data")
	}

	user, err := listingServ.GetUserByUsername(verifying.Username)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(verifying.Password), appConfig.BcryptCost)
	if err != nil {
		return nil, err
	}

	if string(hashedPassword) != user.Hash {
		return nil, errors.New("username or password is not correct")
	}

	bunches, keys, err := listingServ.GetUserBunchKeys(user.ID)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(appConfig.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, tokenClaims{
		bunches,
		keys,
		jwtgo.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth",
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(appConfig.SigningText)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &Token{tokenString}, nil
	}
}

func UpdateUser(ctx context.Context, request interface{}) (interface{}, error) {
	listingServ := ctx.Value(common.ListingServiceContextKey).(listing.Service)
	modifyServ := ctx.Value(common.ModifyServiceContextKey).(modifying.Service)
	appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

	updating, ok := request.(UpdatingUser)
	if !ok {
		return nil, errors.New("couldn't get updating data")
	}

	modifyServ.ModifyUser

}
