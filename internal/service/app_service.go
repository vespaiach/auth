package service

import (
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

// AppService contains all app's services
type AppService struct {
	UserService       *UserService
	TokenService      TokenService
	ActionService     *ActionService
	RoleService       *RoleService
	RoleActionService *RoleActionService
	UserRoleService   *UserRoleService
}

// NewAppService init app service instance
func NewAppService(appRepo *model.AppRepo, appConfig *conf.AppConfig) *AppService {
	return &AppService{
		UserService:       NewUserService(appRepo, appConfig),
		TokenService:      NewTokenService(appRepo, appConfig),
		ActionService:     NewActionService(appRepo, appConfig),
		RoleService:       NewRoleService(appRepo, appConfig),
		RoleActionService: NewRoleActionService(appRepo, appConfig),
		UserRoleService:   NewUserRoleService(appRepo, appConfig),
	}
}
