package store

import "time"

// User model
type User struct {
	ID        uint      `json:"id"`
	Name      string    `gorm:"type:varchar(255)"`
	Username  string    `gorm:"type:varchar(63);unique_index"`
	Hashed    string    `gorm:"type:varchar(511)"`
	Email     string    `gorm:"type:varchar(255);unique_index"`
	Active    int       `gorm:"type:int,index"`
	Verified  bool      `gorm:"type:bool,index"`
	Actions   []*Action `gorm:"many2many:user_actions;"`
	Roles     []*Role   `gorm:"many2many:user_roles;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (sto *mysqlStore) GetUserByUsername(username string) (user *User, err error) {
	user = new(User)
	sto.db.Where("username=?", username).First(&user)

	if user == nil {
		err = ErrDataNotFound
	}

	return
}

func (sto *mysqlStore) GetUserByEmail(email string) (user *User, err error) {
	user = new(User)
	sto.db.Where("email=?", email).First(&user)

	if user == nil {
		err = ErrDataNotFound
	}

	return
}

func (sto *mysqlStore) CreateUser(name string, username string, hashed string, email string) (*User, error) {
	user := User{
		Name:     name,
		Username: username,
		Hashed:   hashed,
		Email:    email,
	}
	sto.db.Create(&user)

	if sto.db.NewRecord(user) {
		return nil, ErrCanNotCreateData
	}

	return &user, nil
}
