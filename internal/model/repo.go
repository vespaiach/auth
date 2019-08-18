package model

// AppRepo will holds all available repos
type AppRepo struct {
	ActionRepo       ActionRepo
	RoleRepo         RoleRepo
	UserRepo         UserRepo
	UserActionRepo   UserActionRepo
	UserRoleRepo     UserRoleRepo
	RoleActionRepo   RoleActionRepo
	TokenHistoryRepo TokenHistoryRepo
}
