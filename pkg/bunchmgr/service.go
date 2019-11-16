package bunchmgr

import (
	"database/sql"
	"github.com/vespaiach/auth/pkg/common"
	"regexp"
	"strings"
)

type Storer interface {
	AddBunch(name string, desc string) (int64, error)
	ModifyBunch(id int64, name string, desc string, active sql.NullBool) error
	GetBunchByName(name string) (*Bunch, error)
	GetBunch(id int64) (*Bunch, error)
	QueryBunches(take int64, skip int64, name string, active sql.NullBool, sortby string,
		direction common.SortingDirection) ([]*Bunch, int64, error)
	GetKeyIDs(keys []string) ([]int64, error)
	AddKeysToBunch(bunchID int64, keyIDs []int64) error
	GetKeysInBunch(name string) ([]*Key, error)
}

type Service interface {
	AddBunch(name string, desc string) (int64, error)
	ModifyBunch(id int64, name string, desc string, active sql.NullBool) error
	GetBunchByName(name string) (*Bunch, error)
	GetBunch(id int64) (*Bunch, error)
	GetKeysInBunch(name string) ([]*Key, error)
	QueryBunches(page int64, perPage int64, name string, active sql.NullBool, order string) ([]*Bunch, int64, error)
	AddKeysToBunch(bunch string, keys []string) error
}

type service struct {
	st Storer
}

func NewService(st Storer) Service {
	return &service{st}
}

func (s *service) AddBunch(name string, desc string) (int64, error) {
	if !s.isValidKey(name) {
		return 0, common.ErrBunchNameInvalid
	}

	dup, err := s.isDuplicatedKey(name)
	if err != nil {
		return 0, err
	}
	if dup {
		return 0, common.ErrDuplicatedBunch
	}

	return s.st.AddBunch(name, desc)
}

func (s *service) ModifyBunch(id int64, name string, desc string, active sql.NullBool) error {
	updating, err := s.st.GetBunch(id)
	if err != nil {
		return err
	}
	if updating == nil {
		return common.ErrBunchNotFound
	}

	if len(name) == 0 {
		name = updating.Name
	} else {
		if !s.isValidKey(name) {
			return common.ErrBunchNameInvalid
		}
	}

	if len(desc) == 0 {
		desc = updating.Desc
	}

	if updating.Name != name {
		dup, err := s.isDuplicatedKey(name)
		if err != nil {
			return err
		}
		if dup {
			return common.ErrDuplicatedBunch
		}
	}

	return s.st.ModifyBunch(id, name, desc, active)
}

func (s *service) GetBunchByName(name string) (*Bunch, error) {
	return s.st.GetBunchByName(name)
}

func (s *service) GetBunch(id int64) (*Bunch, error) {
	return s.st.GetBunch(id)
}

func (s *service) QueryBunches(page int64, perPage int64, name string, active sql.NullBool,
	order string) ([]*Bunch, int64, error) {

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

	return s.st.QueryBunches(take, skip, name, active, sortby, direction)
}

func (s *service) AddKeysToBunch(bunchName string, keys []string) error {
	bunch, err := s.st.GetBunchByName(bunchName)
	if err != nil {
		return err
	}
	if bunch == nil {
		return common.ErrBunchNotFound
	}

	if len(keys) > 0 {
		keyIDs, err := s.st.GetKeyIDs(keys)
		if err != nil {
			return err
		}
		if len(keyIDs) == 0 {
			return common.ErrKeyNotFound
		}

		return s.st.AddKeysToBunch(bunch.ID, keyIDs)
	}

	return nil
}

func (s *service) GetKeysInBunch(name string) ([]*Key, error) {
	return s.st.GetKeysInBunch(name)
}

func (s *service) isDuplicatedKey(name string) (bool, error) {
	existing, err := s.st.GetBunchByName(name)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

func (s *service) isValidKey(name string) bool {
	matched, err := regexp.Match(`^[a-z0-9_]{1,32}$`, []byte(name))
	return err == nil && matched
}
