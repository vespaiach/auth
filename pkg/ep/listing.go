package ep

import (
	"context"
	"errors"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/listing"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(common.ListingServiceContextKey).(listing.Service)

	id, ok := request.(int64)
	if !ok {
		return nil, errors.New("couldn't get id")
	}

	user, err := service.GetUser(id)
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
