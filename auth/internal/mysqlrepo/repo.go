package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/models"
)

// MysqlAppRepo return all repos implemented by mysql
type MysqlAppRepo struct {
	UserRepo models.UserRepo
}

// NewMysqlAppRepo inits all repos
func NewMysqlAppRepo(db *gorm.DB) *MysqlAppRepo {
	return &MysqlAppRepo{
		UserRepo: NewMysqlUserRepo(db),
	}
}
