package mysqlrepo

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

var userFilterKeys = map[string]bool{"full_name": true, "username": true, "email": true}
var userSortKeys = map[string]bool{"full_name": true, "username": true, "email": true, "created_at": true, "updated_at": true}

// MysqlUserRepo will implement model.UserRepo
type MysqlUserRepo struct {
	DbClient *gorm.DB
}

// NewMysqlUserRepo create new instance of MysqlUserRepo
func NewMysqlUserRepo(db *gorm.DB) model.UserRepo {
	return &MysqlUserRepo{
		db,
	}
}

// GetByID find a user by its ID
func (repo *MysqlUserRepo) GetByID(id uint) (user *model.User, err error) {
	user = new(model.User)
	repo.DbClient.Where("id=?", id).First(&user)

	if user == nil {
		err = comtype.ErrDataNotFound
	}

	return
}

// GetByUsername find a user by its username
func (repo *MysqlUserRepo) GetByUsername(username string) (user *model.User, err error) {
	user = new(model.User)
	repo.DbClient.Where("username=?", username).First(&user)

	if user == nil {
		err = comtype.ErrDataNotFound
	}

	return
}

// GetByEmail find a user by its email
func (repo *MysqlUserRepo) GetByEmail(email string) (user *model.User, err error) {
	user = new(model.User)
	repo.DbClient.Where("email=?", email).First(&user)

	if user == nil {
		err = comtype.ErrDataNotFound
	}

	return
}

// Create a new user
func (repo *MysqlUserRepo) Create(fullName string, username string, hashed string, email string) (*model.User, error) {
	user := model.User{
		FullName: fullName,
		Username: username,
		Hashed:   hashed,
		Email:    email,
	}

	repo.DbClient.Create(&user)

	if repo.DbClient.NewRecord(user) {
		return nil, comtype.ErrCreadDataFailed
	}

	return &user, nil
}

// Update user
func (repo *MysqlUserRepo) Update(id uint, fields map[string]interface{}) (err error) {
	user, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	fields["updated_at"] = time.Now()

	repo.DbClient.Model(&user).Updates(fields)

	return nil
}

// Query a list of users
func (repo *MysqlUserRepo) Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.User, int64, error) {
	db := repo.DbClient.Model(&model.User{})

	var total int64
	users := []*model.User{}
	offset := perPage * (page - 1)

	for k, v := range filters {
		_, ok := userFilterKeys[k]
		if ok {
			if k == "active" {
				if v == comtype.Active {
					db = db.Where("active = ?", 1)
				} else {
					db = db.Where("active = ?", 0)
				}
			} else {
				s, good := v.(string)
				if good {
					db = db.Where(k+" LIKE ?", s+"%")
				} else {
					return nil, 0, comtype.ErrDataTypeMismatch
				}
			}
		} else {
			return nil, 0, comtype.ErrNotAllowField
		}
	}

	for k, v := range sorts {
		_, ok := userSortKeys[k]
		if ok {
			if v == comtype.Ascending {
				db = db.Order(k + " asc")
			} else {
				db = db.Order(k + " desc")
			}
		} else {
			return nil, 0, comtype.ErrNotAllowField
		}
	}

	db.
		Offset(offset).
		Limit(perPage).
		Count(&total).
		Find(&users)

	return users, total, nil
}
