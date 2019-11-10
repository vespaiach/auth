package modifying

// Repository defines modifying's functions
type Repository interface {
	ModifyServiceKey(key ServiceKey) error
	GetKeyByID(id int64) (string, error)
	IsDuplicatedKey(key string) bool
	ModifyBunch(b Bunch) error
	GetBunchNameByID(id int64) (string, error)
	IsDuplicatedBunch(name string) bool
	ModifyUser(u User) error
	// GetUser gets user's username and user's email buy user's id
	GetUsernameAndEmail(id int64) (string, string, error)
	IsDuplicatedUsername(name string) bool
	IsDuplicatedEmail(email string) bool
}

// Service provides modifying operations.
type Service interface {
	ModifyServiceKey(sk ServiceKey) error
	ModifyBunch(b Bunch) error
	ModifyUser(u User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) ModifyServiceKey(sk ServiceKey) error {
	if err := sk.Validate(); err != nil {
		return err
	}

	key, err := s.repo.GetKeyByID(sk.ID)
	if err != nil {
		return err
	}

	if sk.Key != key && s.repo.IsDuplicatedKey(sk.Key) {
		return ErrDuplicatedKey
	}

	return s.repo.ModifyServiceKey(sk)
}

func (s *service) ModifyBunch(b Bunch) error {
	if err := b.Validate(); err != nil {
		return err
	}

	name, err := s.repo.GetBunchNameByID(b.ID)
	if err != nil {
		return err
	}

	if name != b.Name && s.repo.IsDuplicatedBunch(b.Name) {
		return ErrDuplicatedBunch
	}

	return s.repo.ModifyBunch(b)
}

func (s *service) ModifyUser(u User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	username, email, err := s.repo.GetUsernameAndEmail(u.ID)
	if err != nil {
		return err
	}

	if u.Username != username && s.repo.IsDuplicatedUsername(u.Username) {
		return ErrDuplicatedUsername
	}

	if u.Email != email && s.repo.IsDuplicatedUsername(u.Email) {
		return ErrDuplicatedEmail
	}

	return s.repo.ModifyUser(u)
}
