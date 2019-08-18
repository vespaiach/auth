package mysqlrepo

import (
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/internal/model"
)

// NewMysqlAppRepo inits all repos with mysql
func NewMysqlAppRepo(db *sqlx.DB) *model.AppRepo {
	return &model.AppRepo{
		UserRepo:         NewMysqlUserRepo(db),
		ActionRepo:       NewMysqlActionRepo(db),
		RoleRepo:         NewMysqlRoleRepo(db),
		UserActionRepo:   NewMysqlUserActionRepo(db),
		UserRoleRepo:     NewMysqlUserRoleRepo(db),
		RoleActionRepo:   NewMysqlRoleActionRepo(db),
		TokenHistoryRepo: NewMysqlTokenHistoryRepo(db),
	}
}
