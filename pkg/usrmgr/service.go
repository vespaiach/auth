package usrmgr

import (
	"database/sql"
	"github.com/vespaiach/auth/pkg/common"
	"regexp"
	"strings"
)

var emailReg = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Storer interface {
	AddUser(username string, email string, hash string) (int64, error)
	ModifyUser(id int64, username string, email string, hash string, active sql.NullBool) error
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUser(id int64) (*User, error)
	AddBunchesToUser(userID int64, bunchIDs []int64) error
	QueryUsers(take int64, skip int64, username string, email string, active sql.NullBool, sortby string,
		direction common.SortingDirection) ([]*User, int64, error)
	GetBunchIDs(bunches []string) ([]int64, error)
	GetBunches(username string) ([]*Bunch, error)
	GetKeys(username string) ([]*Key, error)
}

type Service interface {
	AddUser(username string, email string, hash string) (int64, error)
	ModifyUser(id int64, username string, email string, hash string, active sql.NullBool) error
	GetUserByUsername(username string) (*User, error)
	GetUser(id int64) (*User, error)
	QueryUsers(page int64, perPage int64, username string, email string, active sql.NullBool,
		order string) ([]*User, int64, error)
	AddBunchesToUser(username string, bunches []string) error
	GetBunches(username string) ([]*Bunch, error)
	GetKeys(username string) ([]*Key, error)
}

type service struct {
	st Storer
}

func NewService(st Storer) Service {
	return &service{st}
}

func (s *service) AddUser(username string, email string, hash string) (int64, error) {
	if !s.isValidName(username) {
		return 0, common.ErrUsernameInvalid
	}

	if !s.isValidEmail(email) {
		return 0, common.ErrEmailInvalid
	}

	dupName, err := s.isDuplicatedUsername(username)
	if err != nil {
		return 0, err
	}
	if dupName {
		return 0, common.ErrDuplicatedUsername
	}

	dupEmail, err := s.isDuplicatedEmail(email)
	if err != nil {
		return 0, err
	}
	if dupEmail {
		return 0, common.ErrDuplicatedEmail
	}

	if len(hash) == 0 {
		return 0, common.ErrMissingHash
	}

	return s.st.AddUser(username, email, hash)
}

func (s *service) ModifyUser(id int64, username string, email string, hash string, active sql.NullBool) error {
	updating, err := s.st.GetUser(id)
	if err != nil {
		return err
	}
	if updating == nil {
		return common.ErrUserNotFound
	}

	if len(username) > 0 {
		if !s.isValidName(username) {
			return common.ErrUsernameInvalid
		}
		if username != updating.Username {
			dup, err := s.isDuplicatedUsername(username)
			if err != nil {
				return err
			}
			if dup {
				return common.ErrDuplicatedUsername
			}
		}
	}

	if len(email) > 0 {
		if !s.isValidEmail(email) {
			return common.ErrEmailInvalid
		}
		if email != updating.Email {
			dup, err := s.isDuplicatedEmail(email)
			if err != nil {
				return err
			}
			if dup {
				return common.ErrDuplicatedEmail
			}
		}
	}

	return s.st.ModifyUser(id, username, email, hash, active)
}

func (s *service) GetUserByUsername(username string) (*User, error) {
	return s.st.GetUserByUsername(username)
}

func (s *service) GetUser(id int64) (*User, error) {
	return s.st.GetUser(id)
}

func (s *service) QueryUsers(page int64, perPage int64, username string, email string, active sql.NullBool,
	order string) ([]*User, int64, error) {

	var (
		sortby    string
		direction common.SortingDirection
		take      int64 = perPage
		skip      int64 = perPage * (page - 1)
	)

	if len(order) > 0 {
		switch order[0] {
		case '+':
			direction = common.Ascending
			sortby = strings.TrimSpace(order[1:])
			break
		case '-':
			direction = common.Descending
			sortby = strings.TrimSpace(order[1:])
			break
		default:
			direction = common.Descending
			sortby = strings.TrimSpace(order)
			break
		}
	} else {
		sortby = "created_at"
		direction = common.Descending
	}

	return s.st.QueryUsers(take, skip, username, email, active, sortby, direction)
}

func (s *service) AddBunchesToUser(username string, bunches []string) error {
	user, err := s.st.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if user == nil {
		return common.ErrUserNotFound
	}

	if len(bunches) > 0 {
		bunchIDs, err := s.st.GetBunchIDs(bunches)
		if err != nil {
			return err
		}
		if len(bunchIDs) == 0 {
			return common.ErrKeyNotFound
		}

		return s.st.AddBunchesToUser(user.ID, bunchIDs)
	}

	return nil
}

func (s *service) GetBunches(username string) ([]*Bunch, error) {
	return s.st.GetBunches(username)
}

func (s *service) GetKeys(username string) ([]*Key, error) {
	return s.st.GetKeys(username)
}

func (s *service) isDuplicatedUsername(username string) (bool, error) {
	existing, err := s.st.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

func (s *service) isDuplicatedEmail(email string) (bool, error) {
	existing, err := s.st.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

func (s *service) isValidName(name string) bool {
	matched, err := regexp.Match(`^[a-zA-Z0-9_]{1,32}$`, []byte(name))
	return err == nil && matched
}

func (s *service) isValidEmail(email string) bool {
	if !emailReg.MatchString(email) {
		return false
	}
	return true
}
