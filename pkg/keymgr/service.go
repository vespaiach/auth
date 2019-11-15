package keymgr

import "regexp"

type Storer interface {
	AddKey(name string, desc string) (int64, error)
	GetKeyByName(name string) (*Key, error)
	GetKey(id int64) (*Key, error)
	GetBunchID(name string) (int64, error)
	ModifyKey(id int64, name string) (success bool, err error)
	AddKeyToBunch(keyID int64, bunchID int64) (int64, error)
}

type Service interface {
	AddKey(name string, desc string) (id int64, err error)
	ModifyKey(id int64, name string) (success bool, err error)
	GetKey(id int64) (key *Key, err error)
	GetKeyByName(name string) (key *Key, err error)
	AddKeyToBunch(name string, bunch string) (int64, error)
}

type service struct {
	st Storer
}

func NewService(st Storer) Service {
	return &service{st}
}

func (s *service) AddKey(name string, desc string) (int64, error) {
	if !s.isValidKey(name) {
		return 0, ErrKeyNameInvalid
	}

	existing, err := s.st.GetKeyByName(name)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, ErrDuplicatedKey
	}

	return s.st.AddKey(name, desc)
}

func (s *service) ModifyKey(id int64, name string) (bool, error) {
	if !s.isValidKey(name) {
		return false, ErrKeyNameInvalid
	}

	updating, err := s.st.GetKey(id)
	if err != nil {
		return false, err
	}
	if updating == nil {
		return false, ErrKeyNotFound
	}

	if updating.Key == name {
		return true, nil
	}

	dup, err := s.isDuplicatedKey(name)
	if err != nil {
		return false, err
	}
	if dup {
		return false, ErrDuplicatedKey
	}

	return s.st.ModifyKey(id, name)
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
		return 0, ErrKeyNotFound
	}

	bunchID, err := s.st.GetBunchID(name)
	if err != nil {
		return 0, err
	}
	if bunchID == 0 {
		return 0, ErrBunchNotFound
	}

	return s.st.AddKeyToBunch(key.ID, bunchID)
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
