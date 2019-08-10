package mysqlrepo

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

var actionFilterKeys = map[string]bool{"name": true, "active": true}
var actionSortKeys = map[string]bool{"name": true, "active": true, "created_at": true}

// MysqlActionRepo will implement model.ActionRepo
type MysqlActionRepo struct {
	DbClient *gorm.DB
}

// NewMysqlActionRepo create new instance of MysqlActionRepo
func NewMysqlActionRepo(db *gorm.DB) model.ActionRepo {
	return &MysqlActionRepo{
		db,
	}
}

// GetByID find an action by its ID
func (repo *MysqlActionRepo) GetByID(id uint) (action *model.Action, err error) {
	repo.DbClient.Where("id=?", id).First(action)

	if action == nil {
		err = comtype.ErrDataNotFound
	}

	return
}

// GetByName find an action by its name
func (repo *MysqlActionRepo) GetByName(name string) (action *model.Action, err error) {
	repo.DbClient.Where("username=?", name).First(action)

	if action == nil {
		err = comtype.ErrDataNotFound
	}

	return
}

// Create an new action
func (repo *MysqlActionRepo) Create(name string, desc string) (action *model.Action, err error) {
	action = &model.Action{
		ActionName: name,
		ActionDesc: desc,
	}

	repo.DbClient.Create(action)

	if repo.DbClient.NewRecord(action) {
		action = nil
		err = comtype.ErrCreadDataFailed
	}

	return
}

// Update action
func (repo *MysqlActionRepo) Update(id uint, fields map[string]interface{}) (err error) {
	action, err := repo.GetByID(id)

	if err == nil {
		fields["updated_at"] = time.Now()

		repo.DbClient.Model(action).Updates(fields)
	}

	return
}

// Query a list of actions
func (repo *MysqlActionRepo) Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.Action, int64, error) {
	db := repo.DbClient.Model(&model.User{})

	var total int64
	actions := []*model.Action{}
	offset := perPage * (page - 1)

	for k, v := range filters {
		_, ok := actionFilterKeys[k]
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
		Find(&actions)

	return actions, total, nil
}
