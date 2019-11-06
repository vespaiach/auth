package endpoint

import (
	"context"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/service"
)

// MakeCreateRoleEndpoint is to create a creating-role endpoint
func MakeCreateRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(CreateRoleRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeCreateRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	role, err := service.RoleService.CreateRole(req.RoleName, req.RoleDesc)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeCreateRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return RoleResponse{
			role.ID,
			role.RoleName,
			role.RoleDesc,
			role.Active,
			role.CreatedAt,
			role.UpdatedAt,
		}, nil
	}
}

// MakeUpdateRoleEndpoint is to create a updating-role endpoint
func MakeUpdateRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(UpdateRoleRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeUpdateRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	role, err := service.RoleService.UpdateRole(req.ID, req.RoleName, req.RoleDesc, req.Active)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeUpdateRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return RoleResponse{
			role.ID,
			role.RoleName,
			role.RoleDesc,
			role.Active,
			role.CreatedAt,
			role.UpdatedAt,
		}, nil
	}
}

// MakeGetRoleEndpoint is to create a get-role-by-id endpoint
func MakeGetRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	id, ok := request.(int64)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeGetRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	role, err := service.RoleService.GetRole(id)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeGetRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return RoleResponse{
			role.ID,
			role.RoleName,
			role.RoleDesc,
			role.Active,
			role.CreatedAt,
			role.UpdatedAt,
		}, nil
	}
}

// MakeQueryRoleEndpoint is to create a query-role endpoint
func MakeQueryRoleEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(QueryRoleRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeQueryRoleEndpoint:", comtype.ErrBadRequest, nil)
	}

	roles, err := service.RoleService.FetchRoles(req.Take, req.RoleName, req.Active, req.SortBy)
	if err != nil {
		return nil, err
	}

	res := make([]RoleResponse, 0, len(roles))
	for _, ac := range roles {
		res = append(res, RoleResponse{
			ac.ID,
			ac.RoleName,
			ac.RoleDesc,
			ac.Active,
			ac.CreatedAt,
			ac.UpdatedAt,
		})
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeQueryRoleEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return res, nil
	}
}
