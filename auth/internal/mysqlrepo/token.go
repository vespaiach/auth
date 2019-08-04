package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlTokenRepo will implement model.TokenRepo
type MysqlTokenRepo struct {
	DbClient *gorm.DB
}

// NewMysqlTokenRepo create new instance of MysqlTokenRepo
func NewMysqlTokenRepo(db *gorm.DB) model.TokenRepo {
	return &MysqlTokenRepo{
		db,
	}
}

// GetByID find a user by its ID
func (repo *MysqlTokenRepo) GetByID(id string) (token *model.Token, err error) {
	token = new(model.Token)
	repo.DbClient.Where("id=?", id).First(&token)

	if token == nil {
		err = comtype.ErrDataNotFound
	}

	return
}

// Save a token
func (repo *MysqlTokenRepo) Save(id string, userID uint, accessToken string, refreshToken string) error {
	token := model.Token{
		ID:           id,
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	repo.DbClient.Create(&token)

	if repo.DbClient.NewRecord(token) {
		return comtype.ErrCreadDataFailed
	}

	return nil
}
