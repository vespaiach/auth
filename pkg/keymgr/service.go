package keymgr

import (
	"github.com/vespaiach/auth/pkg/common"
	"regexp"
	"strings"
)

type Storer interface {
	AddKey(name string, desc string) (int64, error)
	GetKeyByName(name string) (*Key, error)
	GetKey(id int64) (*Key, error)
	GetBunchID(name string) (int64, error)
	ModifyKey(id int64, name string, desc string) error
	AddKeyToBunch(keyID int64, bunchID int64) (int64, error)
	QueryKeys(take int64, skip int64, name string, sortby string,
		direction common.SortingDirection) ([]*Key, int64, error)
}

type Service interface {
	AddKey(name string, desc string) (id int64, err error)
	ModifyKey(id int64, name string, desc string) (err error)
	GetKey(id int64) (key *Key, err error)
	GetKeyByName(name string) (key *Key, err error)
	AddKeyToBunch(name string, bunch string) (int64, error)
	QueryKeys(page int64, perPage int64, name string, order string) ([]*Key, int64, error)
}

type service struct {
	st Storer
}

func NewService(st Storer) Service {
	return &service{st}
}

func (s *service) AddKey(name string, desc string) (int64, error) {
	if !s.isValidKey(name) {
		return 0, common.ErrKeyNameInvalid
	}

	existing, err := s.st.GetKeyByName(name)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, common.ErrDuplicatedKey
	}

	return s.st.AddKey(name, desc)
}

func (s *service) ModifyKey(id int64, name string, desc string) error {
	updating, err := s.st.GetKey(id)
	if err != nil {
		return err
	}
	if updating == nil {
		return common.ErrKeyNotFound
	}

	if len(name) == 0 {
		name = updating.Key
	} else {
		if !s.isValidKey(name) {
			return common.ErrKeyNameInvalid
		}
	}

	if len(desc) == 0 {
		desc = updating.Desc
	}

	if updating.Key != name {
		dup, err := s.isDuplicatedKey(name)
		if err != nil {
			return err
		}
		if dup {
			return common.ErrDuplicatedKey
		}
	}

	return s.st.ModifyKey(id, name, desc)
}

func (s *service) GetKey(id int64) (*Key, error) {
	return s.st.GetKey(id)
}

func (s *service) GetKeyByName(name string) (*Key, error) {
	return s.st.GetKeyByName(name)
}

func (s *service) AddKeyToBunch(name string, bunch string) (int64, error) {
	key, err := s.st.GetKeyByName(name)
	if err != nil {
		return 0, err
	}
	if key == nil {
		return 0, common.ErrKeyNotFound
	}

	bunchID, err := s.st.GetBunchID(bunch)
	if err != nil {
		return 0, err
	}
	if bunchID == 0 {
		return 0, common.ErrBunchNotFound
	}

	return s.st.AddKeyToBunch(key.ID, bunchID)
}

func (s *service) QueryKeys(page int64, perPage int64, name string, order string) ([]*Key, int64, error) {
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

	return s.st.QueryKeys(take, skip, name, sortby, direction)
}

func (s *service) isDuplicatedKey(name string) (bool, error) {
	existing, err := s.st.GetKeyByName(name)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

func (s *service) isValidKey(name string) bool {
	matched, err := regexp.Match(`^[a-z0-9_]{1,32}$`, []byte(name))
	return err == nil && matched
}
