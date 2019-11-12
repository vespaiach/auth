package adding

import (
	"errors"

	"github.com/vespaiach/auth/pkg/common"
)

// Repository defines storage's functions
type Repository interface {
	AddServiceKey(key ServiceKey) (int64, error)
	IsDuplicatedKey(key string) bool
	AddBunch(b Bunch) (int64, error)
	IsDuplicatedBunch(name string) bool
	AddUser(u User) (int64, error)
	IsDuplicatedUsername(name string) bool
	IsDuplicatedEmail(email string) bool
}

// Service provides service_key, user, role, adding operations.
type Service interface {
	AddServiceKey(key ServiceKey) (int64, error)
	AddBunch(b Bunch) (int64, error)
	AddUser(u User) (int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) AddServiceKey(sk ServiceKey) (int64, error) {
	if err := sk.Validate(); err != nil {
		return 0, err
	}

	if s.repo.IsDuplicatedKey(sk.Key) {
		return 0, common.NewAppErr(errors.New("key is duplicated"), common.ErrDataFailValidation)
	}

	return s.repo.AddServiceKey(sk)
}

func (s *service) AddBunch(b Bunch) (int64, error) {
	if err := b.Validate(); err != nil {
		return 0, err
	}

	if s.repo.IsDuplicatedBunch(b.Name) {
		return 0, common.NewAppErr(errors.New("bunch is duplicated"), common.ErrDataFailValidation)
	}

	return s.repo.AddBunch(b)
}

func (s *service) AddUser(u User) (int64, error) {
	if err := u.Validate(); err != nil {
		return 0, err
	}

	if s.repo.IsDuplicatedUsername(u.Username) {
		return 0, common.NewAppErr(errors.New("username is duplicated"), common.ErrDataFailValidation)
	}

	if s.repo.IsDuplicatedUsername(u.Email) {
		return 0, common.NewAppErr(errors.New("email address is duplicated"), common.ErrDataFailValidation)
	}

	return s.repo.AddUser(u)
}
