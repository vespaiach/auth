package endpoint

import (
	"context"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/service"
)

// MakeCreateRoleActionEndpoint is to create a role-action endpoint
func MakeCreateRoleActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(CreateRoleActionRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeCreateRoleActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	roleAction, err := service.RoleActionService.CreateRoleAction(req.RoleID, req.ActionID)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeCreateRoleActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return RoleActionResponse{
			roleAction.ID,
			roleAction.RoleID,
			roleAction.ActionID,
			RoleResponse{
				roleAction.Role.ID,
				roleAction.Role.RoleName,
				roleAction.Role.RoleDesc,
				roleAction.Role.Active,
				roleAction.Role.CreatedAt,
				roleAction.Role.UpdatedAt,
			},
			ActionResponse{
				roleAction.Action.ID,
				roleAction.Action.ActionName,
				roleAction.Action.ActionDesc,
				roleAction.Action.Active,
				roleAction.Action.CreatedAt,
				roleAction.Action.UpdatedAt,
			},
		}, nil
	}
}

// MakeDeleteRoleActionEndpoint is to create a deleting role-action endpoint
func MakeDeleteRoleActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(DeleteRoleActionRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeDeleteRoleActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	err := service.RoleActionService.DeleteRoleAction(req.RoleID, req.ActionID)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeDeleteRoleActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return true, nil
	}
}

// MakeGetRoleActionEndpoint is to create a get-role-action-by-id endpoint
func MakeGetRoleActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	id, ok := request.(int64)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeGetRoleActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	roleAction, err := service.RoleActionService.GetRoleAction(id)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeGetRoleActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return RoleActionResponse{
			roleAction.ID,
			roleAction.RoleID,
			roleAction.ActionID,
			RoleResponse{
				roleAction.Role.ID,
				roleAction.Role.RoleName,
				roleAction.Role.RoleDesc,
				roleAction.Role.Active,
				roleAction.Role.CreatedAt,
				roleAction.Role.UpdatedAt,
			},
			ActionResponse{
				roleAction.Action.ID,
				roleAction.Action.ActionName,
				roleAction.Action.ActionDesc,
				roleAction.Action.Active,
				roleAction.Action.CreatedAt,
				roleAction.Action.UpdatedAt,
			},
		}, nil
	}
}

// MakeQueryRoleActionEndpoint is to create a query-action endpoint
func MakeQueryRoleActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(QueryRoleActionRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeQueryRoleActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	roleActions, err := service.RoleActionService.FetchRoleActions(req.Take, req.RoleID, req.ActionID)
	if err != nil {
		return nil, err
	}

	res := make([]RoleActionResponse, 0, len(roleActions))
	for _, ac := range roleActions {
		res = append(res, RoleActionResponse{
			ac.ID,
			ac.RoleID,
			ac.ActionID,
			RoleResponse{
				ac.Role.ID,
				ac.Role.RoleName,
				ac.Role.RoleDesc,
				ac.Role.Active,
				ac.Role.CreatedAt,
				ac.Role.UpdatedAt,
			},
			ActionResponse{
				ac.Action.ID,
				ac.Action.ActionName,
				ac.Action.ActionDesc,
				ac.Action.Active,
				ac.Action.CreatedAt,
				ac.Action.UpdatedAt,
			},
		})
	}
	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeQueryRoleActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return res, nil
	}
}
