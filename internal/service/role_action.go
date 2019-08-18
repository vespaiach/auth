package service

import (
	"sync"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

// RoleActionService role's service
type RoleActionService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// NewRoleActionService creates a struct of role-action service
func NewRoleActionService(appRepo *model.AppRepo, appConfig *conf.AppConfig) *RoleActionService {
	return &RoleActionService{
		appConfig,
		appRepo,
	}
}

// CreateRoleAction create a new role-action
func (as *RoleActionService) CreateRoleAction(roleID int64, actionID int64) (*model.RoleAction, *comtype.CommonError) {
	roleActionRepo := as.appRepo.RoleActionRepo
	roleRepo := as.appRepo.RoleRepo
	actionRepo := as.appRepo.ActionRepo

	var (
		wg                 sync.WaitGroup
		role               *model.Role
		action             *model.Action
		actionErr, roleErr *comtype.CommonError
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		role, roleErr = roleRepo.GetByID(roleID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		action, actionErr = actionRepo.GetByID(actionID)
	}()

	wg.Wait()

	if actionErr != nil {
		return nil, actionErr
	}
	if roleErr != nil {
		return nil, roleErr
	}

	errMessages := map[string]string{}
	if role == nil {
		errMessages["role_id"] = "role_id doesn't existing"
	}
	if action == nil {
		errMessages["action_id"] = "action_id doesn't existing"
	}

	if len(errMessages) > 0 {
		return nil, comtype.NewCommonError(nil, "CreateRoleAction", comtype.ErrDataNotFound, errMessages)
	}

	id, err := roleActionRepo.Create(role.ID, action.ID)
	if err != nil {
		return nil, err
	}

	return roleActionRepo.GetByID(id)
}

// FetchRoleActions gets a list of role-actions
func (as *RoleActionService) FetchRoleActions(take int, roleID int64, actionID int64) ([]*model.RoleAction, *comtype.CommonError) {
	repo := as.appRepo.RoleActionRepo
	filters := make(map[string]interface{})

	if take == 0 || take > 100 {
		take = 100
	}

	if roleID > 0 {
		filters["role_id"] = roleID
	}

	if actionID > 0 {
		filters["action_id"] = actionID
	}

	return repo.Query(take, filters)
}

// GetRoleAction gets a role-action by id
func (as *RoleActionService) GetRoleAction(id int64) (*model.RoleAction, *comtype.CommonError) {
	return as.appRepo.RoleActionRepo.GetByID(id)
}

// DeleteRoleAction delete role-action
func (as *RoleActionService) DeleteRoleAction(roleID int64, actionID int64) *comtype.CommonError {
	repo := as.appRepo.RoleActionRepo

	roleActions, err := as.FetchRoleActions(1, roleID, actionID)
	if err != nil {
		return err
	}
	if len(roleActions) == 0 {
		return nil
	}

	return repo.Delete(roleActions[0].ID)
}
