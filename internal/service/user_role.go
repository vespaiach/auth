package service

import (
	"sync"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

// UserRoleService role's service
type UserRoleService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// NewUserRoleService creates a struct of user-role service
func NewUserRoleService(appRepo *model.AppRepo, appConfig *conf.AppConfig) *UserRoleService {
	return &UserRoleService{
		appConfig,
		appRepo,
	}
}

// CreateUserRole create a new role-action
func (as *UserRoleService) CreateUserRole(userID int64, roleID int64) (*model.UserRole, *comtype.CommonError) {
	userRoleRepo := as.appRepo.UserRoleRepo
	roleRepo := as.appRepo.RoleRepo
	userRepo := as.appRepo.UserRepo

	var (
		wg               sync.WaitGroup
		role             *model.Role
		user             *model.User
		userErr, roleErr *comtype.CommonError
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		role, roleErr = roleRepo.GetByID(roleID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		user, userErr = userRepo.GetByID(userID)
	}()

	wg.Wait()

	if userErr != nil {
		return nil, userErr
	}
	if roleErr != nil {
		return nil, roleErr
	}

	errMessages := map[string]string{}
	if role == nil {
		errMessages["role_id"] = "role_id doesn't existing"
	}
	if user == nil {
		errMessages["user_id"] = "user_id doesn't existing"
	}

	if len(errMessages) > 0 {
		return nil, comtype.NewCommonError(nil, "CreateUserRole", comtype.ErrDataNotFound, errMessages)
	}

	id, err := userRoleRepo.Create(user.ID, role.ID)
	if err != nil {
		return nil, err
	}

	return userRoleRepo.GetByID(id)
}

// FetchUserRoles gets a list of user-roles
func (as *UserRoleService) FetchUserRoles(take int, userID int64, roleID int64) ([]*model.UserRole, *comtype.CommonError) {
	repo := as.appRepo.UserRoleRepo
	filters := make(map[string]interface{})

	if take == 0 || take > 100 {
		take = 100
	}

	if userID > 0 {
		filters["user_id"] = userID
	}

	if roleID > 0 {
		filters["role_id"] = roleID
	}

	return repo.Query(take, filters)
}

// GetUserRole gets a user-role by id
func (as *UserRoleService) GetUserRole(id int64) (*model.UserRole, *comtype.CommonError) {
	return as.appRepo.UserRoleRepo.GetByID(id)
}

// DeleteUserRole delete user-role
func (as *UserRoleService) DeleteUserRole(userID int64, roleID int64) *comtype.CommonError {
	repo := as.appRepo.UserRoleRepo

	userRoles, err := as.FetchUserRoles(1, userID, roleID)
	if err != nil {
		return err
	}
	if len(userRoles) == 0 {
		return nil
	}

	return repo.Delete(userRoles[0].ID)
}
