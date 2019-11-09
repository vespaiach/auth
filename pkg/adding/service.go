package adding

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

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) AddServiceKey(sk ServiceKey) (int64, error) {
	if err := sk.Validate(); err != nil {
		return 0, err
	}

	if s.repo.IsDuplicatedKey(sk.Key) {
		return 0, ErrDuplicatedKey
	}

	return s.repo.AddServiceKey(sk)
}

func (s *service) AddBunch(b Bunch) (int64, error) {
	if err := b.Validate(); err != nil {
		return 0, err
	}

	if s.repo.IsDuplicatedBunch(b.Name) {
		return 0, ErrDuplicatedBunch
	}

	return s.repo.AddBunch(b)
}

func (s *service) AddUser(u User) (int64, error) {
	if err := u.Validate(); err != nil {
		return 0, err
	}

	if s.repo.IsDuplicatedUsername(u.Username) {
		return 0, ErrDuplicatedUsername
	}

	if s.repo.IsDuplicatedUsername(u.Email) {
		return 0, ErrDuplicatedEmail
	}

	return s.repo.AddUser(u)
}
