package model

// AppRepo constains all available repos
type AppRepo struct {
	ActionRepo       ActionRepo
	RoleRepo         RoleRepo
	UserRepo         UserRepo
	UserActionRepo   UserActionRepo
	UserRoleRepo     UserRoleRepo
	RoleActionRepo   RoleActionRepo
	TokenHistoryRepo TokenHistoryRepo
}
