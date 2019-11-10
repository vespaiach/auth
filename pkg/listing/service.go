package listing

import (
	"errors"
	"github.com/vespaiach/auth/pkg/common"
)

type Repository interface {
	GetUserByID(id int64) (*User, error)
}

type Service interface {
	GetUser(id int64) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetUser(id int64) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, common.NewAppErr(err, common.ErrGetData)
	}

	if user == nil {
		return nil, common.NewAppErr(errors.New("data not found"), common.ErrDataNotFound)
	}

	return user, nil
}
