package endpoint

import (
	"context"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/service"
)

// MakeCreateUserRoleEndpoint is to create a role-action endpoint
func MakeCreateUserRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(CreateUserRoleRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeCreateUserRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	userRole, err := service.UserRoleService.CreateUserRole(req.UserID, req.RoleID)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeCreateUserRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return UserRoleResponse{
			userRole.ID,
			userRole.UserID,
			userRole.RoleID,
			UserResponse{
				userRole.User.ID,
				userRole.User.FullName,
				userRole.User.Username,
				userRole.User.Email,
				userRole.User.Verified,
				userRole.User.Active,
				userRole.User.CreatedAt,
				userRole.User.UpdatedAt,
			},
			RoleResponse{
				userRole.Role.ID,
				userRole.Role.RoleName,
				userRole.Role.RoleDesc,
				userRole.Role.Active,
				userRole.Role.CreatedAt,
				userRole.Role.UpdatedAt,
			},
		}, nil
	}
}

// MakeDeleteUserRoleEndpoint is to create a deleting role-action endpoint
func MakeDeleteUserRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(DeleteUserRoleRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeDeleteUserRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	err := service.UserRoleService.DeleteUserRole(req.UserID, req.RoleID)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeDeleteUserRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return true, nil
	}
}

// MakeGetUserRoleEndpoint is to create a get-role-action-by-id endpoint
func MakeGetUserRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	id, ok := request.(int64)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeGetUserRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	userRole, err := service.UserRoleService.GetUserRole(id)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeGetUserRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return UserRoleResponse{
			userRole.ID,
			userRole.UserID,
			userRole.RoleID,
			UserResponse{
				userRole.User.ID,
				userRole.User.FullName,
				userRole.User.Username,
				userRole.User.Email,
				userRole.User.Verified,
				userRole.User.Active,
				userRole.User.CreatedAt,
				userRole.User.UpdatedAt,
			},
			RoleResponse{
				userRole.Role.ID,
				userRole.Role.RoleName,
				userRole.Role.RoleDesc,
				userRole.Role.Active,
				userRole.Role.CreatedAt,
				userRole.Role.UpdatedAt,
			},
		}, nil
	}
}

// MakeQueryUserRoleEndpoint is to create a query-action endpoint
func MakeQueryUserRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(QueryUserRoleRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeQueryUserRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	userRoles, err := service.UserRoleService.FetchUserRoles(req.Take, req.UserID, req.RoleID)
	if err != nil {
		return nil, err
	}

	res := make([]UserRoleResponse, 0, len(userRoles))
	for _, ur := range userRoles {
		res = append(res, UserRoleResponse{
			ur.ID,
			ur.UserID,
			ur.RoleID,
			UserResponse{
				ur.User.ID,
				ur.User.FullName,
				ur.User.Username,
				ur.User.Email,
				ur.User.Verified,
				ur.User.Active,
				ur.User.CreatedAt,
				ur.User.UpdatedAt,
			},
			RoleResponse{
				ur.Role.ID,
				ur.Role.RoleName,
				ur.Role.RoleDesc,
				ur.Role.Active,
				ur.Role.CreatedAt,
				ur.Role.UpdatedAt,
			},
		})
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeQueryUserRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return res, nil
	}
}
