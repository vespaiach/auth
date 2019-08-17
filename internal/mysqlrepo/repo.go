package mysqlrepo

// import (
// 	"github.com/jinzhu/gorm"
// 	"github.com/vespaiach/auth/internal/model"
// )

// // MysqlAppRepo return all repos implemented by mysql
// type MysqlAppRepo struct {
// 	ActionRepo model.ActionRepo
// }

// // NewMysqlAppRepo inits all repos
// func NewMysqlAppRepo(db *gorm.DB) *model.AppRepo {
// 	return &model.AppRepo{
// 		UserRepo:   NewMysqlUserRepo(db),
// 		ActionRepo: NewMysqlActionRepo(db),
// 		RoleRepo:   NewMysqlRoleRepo(db),
// 	}
// }
