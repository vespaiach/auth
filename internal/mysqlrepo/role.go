package mysqlrepo

// import (
// 	"time"

// 	"github.com/jinzhu/gorm"
// 	"github.com/vespaiach/auth/internal/comtype"
// 	"github.com/vespaiach/auth/internal/model"
// )

// var roleFilterKeys = map[string]bool{"role_name": true, "active": true}
// var roleSortKeys = map[string]bool{"role_name": true, "active": true, "created_at": true}

// // MysqlRoleRepo will implement model.RoleRepo
// type MysqlRoleRepo struct {
// 	DbClient *gorm.DB
// }

// // NewMysqlRoleRepo create new instance of MysqlRoleRepo
// func NewMysqlRoleRepo(db *gorm.DB) model.RoleRepo {
// 	return &MysqlRoleRepo{
// 		db,
// 	}
// }

// // GetByID find a role by its ID
// func (repo *MysqlRoleRepo) GetByID(id int64) (*model.Role, error) {
// 	role := &model.Role{}

// 	if err := repo.DbClient.First(role, "id=?", id).Error; err != nil {
// 		return nil, err
// 	}

// 	role.Actions = []*model.Action{}
// 	repo.DbClient.Model(role).Related(&role.Actions, "Actions")

// 	return role, nil
// }

// // GetByName find a role by its name
// func (repo *MysqlRoleRepo) GetByName(name string) (*model.Role, error) {
// 	role := &model.Role{}

// 	if err := repo.DbClient.First(role, "role_name=?", name).Error; err != nil {
// 		return nil, err
// 	}

// 	role.Actions = []*model.Action{}
// 	repo.DbClient.Model(role).Related(&role.Actions, "Actions")

// 	return role, nil
// }

// // Create a new role
// func (repo *MysqlRoleRepo) Create(name string, desc string) (int64, error) {
// 	role := &model.Role{
// 		RoleName: name,
// 		RoleDesc: desc,
// 	}

// 	if err := repo.DbClient.Create(role).Error; err != nil {
// 		return 0, nil
// 	}

// 	if repo.DbClient.NewRecord(role) {
// 		return 0, comtype.ErrCreadDataFailed
// 	}

// 	return role.ID, nil
// }

// // Update role
// func (repo *MysqlRoleRepo) Update(id int64, fields map[string]interface{}) error {
// 	role, err := repo.GetByID(id)

// 	if err == nil {
// 		fields["updated_at"] = time.Now()
// 		repo.DbClient.Model(role).Updates(fields)

// 		return nil
// 	}

// 	return err
// }

// // Query a list of roles
// func (repo *MysqlRoleRepo) Query(page int64, perPage int64, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.Role, int64, error) {
// 	db := repo.DbClient.Model(&model.Role{})

// 	var total int64
// 	roles := []*model.Role{}
// 	offset := perPage * (page - 1)

// 	for k, v := range filters {
// 		_, ok := roleFilterKeys[k]
// 		if ok {
// 			if k == "active" {
// 				if v == comtype.Active {
// 					db = db.Where("active = ?", 1)
// 				} else {
// 					db = db.Where("active = ?", 0)
// 				}
// 			} else {
// 				s, good := v.(string)
// 				if good {
// 					db = db.Where(k+" LIKE ?", s+"%")
// 				} else {
// 					return nil, 0, comtype.ErrDataTypeMismatch
// 				}
// 			}
// 		} else {
// 			return nil, 0, comtype.ErrNotAllowField
// 		}
// 	}

// 	for k, v := range sorts {
// 		_, ok := userSortKeys[k]
// 		if ok {
// 			if v == comtype.Ascending {
// 				db = db.Order(k + " asc")
// 			} else {
// 				db = db.Order(k + " desc")
// 			}
// 		} else {
// 			return nil, 0, comtype.ErrNotAllowField
// 		}
// 	}

// 	db.
// 		Preload("Actions").
// 		Offset(offset).
// 		Limit(perPage).
// 		Count(&total).
// 		Find(&roles)

// 	return roles, total, nil
// }
