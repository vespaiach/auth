package service

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

// ActionService action's service
type ActionService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// NewActionService creates a struct that implement ActionService interface
func NewActionService(appRepo *model.AppRepo, appConfig *conf.AppConfig) *ActionService {
	return &ActionService{
		appConfig,
		appRepo,
	}
}

// CreateAction create a new action
func (as *ActionService) CreateAction(actionName string, actionDesc string) (*model.Action, *comtype.CommonError) {
	actionRepo := as.appRepo.ActionRepo

	existing, _ := actionRepo.GetByName(actionName)
	if existing != nil {
		return nil, comtype.NewCommonError(nil, "ActionService - CreateAction", comtype.ErrDuplicatedData, map[string]string{"action_name": "action_name is duplicated"})
	}

	id, err := actionRepo.Create(actionName, actionDesc)
	if err != nil {
		return nil, err
	}

	return actionRepo.GetByID(id)
}

// UpdateAction update action
func (as *ActionService) UpdateAction(id int64, actionName string, actionDesc string, active *bool) (*model.Action, *comtype.CommonError) {
	actionRepo := as.appRepo.ActionRepo
	updatingMap := make(map[string]interface{})

	if len(actionName) > 0 {
		updatingMap["action_name"] = actionName

		existing, _ := actionRepo.GetByName(actionName)
		if existing != nil {
			return nil, comtype.NewCommonError(nil, "ActionService - CreateAction",
				comtype.ErrDuplicatedData, map[string]string{"action_name": "action_name is duplicated"})
		}
	}

	if len(actionDesc) > 0 {
		updatingMap["action_desc"] = actionDesc
	}

	if active != nil {
		updatingMap["active"] = *active
	}

	updatingMap["updated_at"] = time.Now()

	err := actionRepo.Update(id, updatingMap)
	if err != nil {
		return nil, err
	}

	return actionRepo.GetByID(id)
}

// GetAction gets an action by ID
func (as *ActionService) GetAction(id int64) (*model.Action, *comtype.CommonError) {
	actionRepo := as.appRepo.ActionRepo
	return actionRepo.GetByID(id)
}

// FetchActions gets a list of actions
func (as *ActionService) FetchActions(take int, actionName string, active *bool, sortBy string) ([]*model.Action, *comtype.CommonError) {
	actionRepo := as.appRepo.ActionRepo
	filters := make(map[string]interface{})
	sorts := make(map[string]comtype.SortDirection)

	if take == 0 || take > 100 {
		take = 100
	}

	if len(actionName) > 0 {
		filters["action_name"] = actionName
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

	return actionRepo.Query(take, filters, sorts)
}
