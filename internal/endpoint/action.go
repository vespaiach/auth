package endpoint

import (
	"context"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/service"
)

// MakeCreateActionEndpoint is to create a creating-action endpoint
func MakeCreateActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)
	req, ok := request.(CreateActionRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeCreateActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	action, err := service.ActionService.CreateAction(req.ActionName, req.ActionDesc)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeCreateActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return ActionResponse{
			action.ID,
			action.ActionName,
			action.ActionDesc,
			action.Active,
			action.CreatedAt,
			action.UpdatedAt,
		}, nil
	}
}

// MakeUpdateActionEndpoint is to create a updating-action endpoint
func MakeUpdateActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)
	req, ok := request.(UpdateActionRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeUpdateActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	action, err := service.ActionService.UpdateAction(req.ID, req.ActionName, req.ActionDesc, req.Active)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeUpdateActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return ActionResponse{
			action.ID,
			action.ActionName,
			action.ActionDesc,
			action.Active,
			action.CreatedAt,
			action.UpdatedAt,
		}, nil
	}
}

// MakeGetActionEndpoint is to create a get-action-by-id endpoint
func MakeGetActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)
	id, ok := request.(int64)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeGetActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	action, err := service.ActionService.GetAction(id)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeGetActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return ActionResponse{
			action.ID,
			action.ActionName,
			action.ActionDesc,
			action.Active,
			action.CreatedAt,
			action.UpdatedAt,
		}, nil
	}
}

// MakeQueryActionEndpoint is to create a query-action endpoint
func MakeQueryActionEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)
	req, ok := request.(QueryActionRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeQueryActionEndpoint:", comtype.ErrBadRequest, nil)
	}

	actions, err := service.ActionService.FetchActions(req.Take, req.ActionName, req.Active, req.SortBy)
	if err != nil {
		return nil, err
	}

	res := make([]ActionResponse, 0, len(actions))
	for _, ac := range actions {
		res = append(res, ActionResponse{
			ac.ID,
			ac.ActionName,
			ac.ActionDesc,
			ac.Active,
			ac.CreatedAt,
			ac.UpdatedAt,
		})
	}
	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeQueryActionEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return res, nil
	}
}
