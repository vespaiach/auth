package store

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/config"
)

// ErrDataNotFound returns an error of data not found
var ErrDataNotFound = errors.New("data not found")

// ErrCanNotCreateData returns an error of not creating data
var ErrCanNotCreateData = errors.New("can not create data")

// Store contains functions to handle data
type Store interface {
	// GetUserByUsername find and return user's data by its username
	GetUserByUsername(username string) (user *User, err error)

	// GetUserByEmail find and return user's data by its email
	GetUserByEmail(email string) (user *User, err error)

	// CreateUser create and return a new user
	CreateUser(name string, username string, hashed string, email string) (user *User, err error)

	// SaveToken will save history of creating token
	SaveToken(id string, userID uint, accessToken string, refreshToken string, action string) error
}

// Action model
type Action struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(63);unique_index"`
	Desc  string  `gorm:"type:varchar(511)"`
	Users []*User `gorm:"many2many:user_actions;"`
	Roles []*Role `gorm:"many2many:role_actions;"`
}

// Role model
type Role struct {
	gorm.Model
	Name    string    `gorm:"type:varchar(63);unique_index"`
	Desc    string    `gorm:"type:varchar(511)"`
	Actions []*Action `gorm:"many2many:role_actions;"`
	Users   []*User   `gorm:"many2many:user_roles;"`
}

type mysqlStore struct {
	config *config.ServiceConfig
	db     *gorm.DB
}

// NewStore get db's connection and create a new store struct
func NewStore(db *gorm.DB, config *config.ServiceConfig) Store {
	return &mysqlStore{
		config,
		db,
	}
}
