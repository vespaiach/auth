package endpoint

import (
	"context"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/service"
)

// MakeRegisterUserEndpoint is to create register-user endpoint
func MakeRegisterUserEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(RegisterUserRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeRegisterUserEndpoint:", comtype.ErrBadRequest, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	user, err := service.UserService.RegisterUser(req.FullName, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeRegisterUserEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return UserResponse{
			user.ID,
			user.FullName,
			user.Username,
			user.Email,
			user.Verified,
			user.Active,
			user.CreatedAt,
			user.UpdatedAt,
		}, nil
	}
}

// MakeUpdateUserEndpoint is to create a updating-user endpoint
func MakeUpdateUserEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(UpdateUserRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeUpdateUserEndpoint:", comtype.ErrDataValidationFail, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	user, err := service.UserService.UpdateUser(req.ID, req.FullName, "", "", "", nil, nil)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeUpdateUserEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return UserResponse{
			user.ID,
			user.FullName,
			user.Username,
			user.Email,
			user.Verified,
			user.Active,
			user.CreatedAt,
			user.UpdatedAt,
		}, nil
	}
}

// MakeChangeUserPasswordEndpoint is to create a updating-user endpoint
func MakeChangeUserPasswordEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	req, ok := request.(ChangeUserPasswordRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeChangeUserPasswordEndpoint:", comtype.ErrDataValidationFail, nil)
	}

	regValidationResults := req.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	user, err := service.UserService.GetUser(req.ID)
	if err != nil {
		return nil, err
	}

	if !service.UserService.IsPasswordMatched(user.Hashed, req.OldPassword) {
		return nil, comtype.NewCommonError(nil, "MakeChangeUserPasswordEndpoint:", comtype.ErrInvalidCredential,
			map[string]string{"old_password": "old_password is not matched"})
	}

	_, err = service.UserService.UpdateUser(req.ID, "", "", "", req.NewPassword, nil, nil)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeChangeUserPasswordEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return true, nil
	}
}

// MakeVerifyUserEndpoint is to create a verifying user's login endpoint
func MakeVerifyUserEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	service := ctx.Value(comtype.CommonKeyRequestContext).(*service.AppService)

	loginReg, ok := request.(VerifyUserRequest)
	if !ok {
		return nil, comtype.NewCommonError(nil, "MakeVerifyUserEndpoint:", comtype.ErrDataValidationFail, nil)
	}

	regValidationResults := loginReg.Validate()
	if regValidationResults != nil {
		return nil, regValidationResults
	}

	user, actions, roles, err := service.UserService.VerifyLogin(loginReg.Username, loginReg.Password)
	if err != nil {
		return nil, err
	}

	token, err := service.TokenService.IssueToken(user, actions, roles, loginReg.RemoteAddr, loginReg.XForwardedFor,
		loginReg.XRealIP, loginReg.UserAgent)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, comtype.NewCommonError(nil, "MakeVerifyUserEndpoint:", comtype.ErrRequestCancelled, nil)
	default:
		return VerifyUserResponse{
			token,
		}, nil
	}
}
