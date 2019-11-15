package listing

import (
	"errors"
	"sync"
)

type Repository interface {
	GetUserByID(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)

	GetBunchByID(id int64) (*Bunch, error)
	GetBunchByName(name string) (*Bunch, error)

	GetKeyByID(id int64) (*Key, error)
	GetKeyByName(key string) (*Key, error)
}

type Service interface {
	GetUser(id int64) (*User, error)
	GetUserBunchKeys(userID int64) ([]string, []string, error)
	GetUserByUsername(name string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

var ErrUserNotFound = errors.New("user is not found")
var ErrGetUserFail = errors.New("couldn't get user data")

func (s *service) GetUser(id int64) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, ErrGetUserFail
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) GetUserByUsername(name string) (*User, error) {
	user, err := s.repo.GetUserByUsername(name)
	if err != nil {
		return nil, ErrGetUserFail
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *service) GetUserBunchKeys(userID int64) ([]string, []string, error) {
	var (
		wg         sync.WaitGroup
		bunches    []*Bunch
		errBunches error
		keys       []*Key
		errKeys    error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bunches, errBunches = s.repo.GetBunchByUserID(userID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		keys, errKeys = s.repo.GetKeysByUserID(userID)
	}()

	wg.Wait()

	if errBunches != nil {
		return nil, nil, errBunches
	}

	if errKeys != nil {
		return nil, nil, errKeys
	}

	bnames := make([]string, 0, len(bunches))
	ks := make([]string, 0, len(keys))

	for _, k := range keys {
		ks = append(ks, k.Key)
	}

	for _, b := range bunches {
		bnames = append(bnames, b.Name)
	}

	return bnames, ks, nil
}
