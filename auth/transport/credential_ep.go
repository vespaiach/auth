package transport

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/vespaiach/authentication/service"
)

type verifyLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type verifyLoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Err          string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

func makeVerifyLoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyLoginRequest)

		u, err := svc.VerifyLogin(req.Username, req.Password)
		if err != nil {
			return verifyLoginResponse{AccessToken: "", RefreshToken: "", Err: err.Error()}, nil
		}

		cred, err := svc.IssueTokens(u)
		if err != nil {
			return verifyLoginResponse{AccessToken: "", RefreshToken: "", Err: err.Error()}, nil
		}

		return verifyLoginResponse{
			AccessToken:  cred.AccessToken,
			RefreshToken: cred.RefreshToken,
			Err:          "",
		}, nil
	}
}

type registerUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type registerUserResponse struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Active    int       `json:"active,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Err       string    `json:"err,omitempty"`
}

func makeRegisterUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(registerUserRequest)

		u, err := svc.RegisterUser(req.Name, req.Username, req.Password, req.Email)
		if err != nil {
			return registerUserResponse{Err: err.Error()}, nil
		}

		return registerUserResponse{
			u.ID, u.Name, u.Username, u.Active, u.Email, u.CreatedAt, u.UpdatedAt, "",
		}, nil
	}
}
