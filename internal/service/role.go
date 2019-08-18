package service

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

// RoleService role's service
type RoleService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// NewRoleService creates a struct of role service
func NewRoleService(appRepo *model.AppRepo, appConfig *conf.AppConfig) *RoleService {
	return &RoleService{
		appConfig,
		appRepo,
	}
}

// CreateRole create a new role
func (as *RoleService) CreateRole(roleName string, roleDesc string) (*model.Role, *comtype.CommonError) {
	repo := as.appRepo.RoleRepo

	existing, _ := repo.GetByName(roleName)
	if existing != nil {
		return nil, comtype.NewCommonError(nil, "RoleService - CreateRole", comtype.ErrDuplicatedData,
			map[string]string{"role_name": "role_name is duplicated"})
	}

	id, err := repo.Create(roleName, roleDesc)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(id)
}

// UpdateRole updates role
func (as *RoleService) UpdateRole(id int64, roleName string, roleDesc string, active *bool) (*model.Role, *comtype.CommonError) {
	repo := as.appRepo.RoleRepo
	updatingMap := make(map[string]interface{})

	if len(roleName) > 0 {
		updatingMap["role_name"] = roleName

		existing, _ := repo.GetByName(roleName)
		if existing != nil {
			return nil, comtype.NewCommonError(nil, "RoleService - UpdateRole",
				comtype.ErrDuplicatedData, map[string]string{"role_name": "role_name is duplicated"})
		}
	}

	if len(roleDesc) > 0 {
		updatingMap["role_desc"] = roleDesc
	}

	if active != nil {
		updatingMap["active"] = *active
	}

	updatingMap["updated_at"] = time.Now()

	err := repo.Update(id, updatingMap)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(id)
}

// GetRole gets an action by ID
func (as *RoleService) GetRole(id int64) (*model.Role, *comtype.CommonError) {
	repo := as.appRepo.RoleRepo
	return repo.GetByID(id)
}

// FetchRoles gets a list of actions
func (as *RoleService) FetchRoles(take int, roleName string, active *bool, sortBy string) ([]*model.Role, *comtype.CommonError) {
	repo := as.appRepo.RoleRepo
	filters := make(map[string]interface{})
	sorts := make(map[string]comtype.SortDirection)

	if take == 0 || take > 100 {
		take = 100
	}

	if len(roleName) > 0 {
		filters["role_desc"] = roleName
	}

	if active != nil {
		filters["active"] = *active
	}

	if len(sortBy) > 0 {
		switch sortBy[0] {
		case '+':
			sorts[sortBy[1:]] = comtype.Ascending
			break
		case '-':
			sorts[sortBy[1:]] = comtype.Decending
			break
		default:
			sorts["created_at"] = comtype.Decending
			break
		}
	} else {
		sorts["created_at"] = comtype.Decending
	}

	return repo.Query(take, filters, sorts)
}
