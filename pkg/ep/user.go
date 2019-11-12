package ep

import (
	"context"
	"errors"

	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/listing"
	"golang.org/x/crypto/bcrypt"
)

// AddingUser model for adding a user
type AddingUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AddUser endpoint is to make an endpoint for adding a new user
func AddUser(ctx context.Context, request interface{}) (interface{}, error) {
	addingServ := ctx.Value(common.AddingServiceContextKey).(adding.Service)
	listingServ := ctx.Value(common.ListingServiceContextKey).(listing.Service)
	appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

	addingUser, ok := request.(AddingUser)
	if !ok {
		return nil, errors.New("couldn't get user data")
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
