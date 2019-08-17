package endpoint

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/shipping/cargo"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/service"
)

type createUserRequest struct {
	FullName string `json:"full_name" validate:"required,lt=64"`
	Username string `json:"username" validate:"required,lt=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=8"`
}

type createUserResponse struct {
	ID        int64     `json:"id,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Error     error     `json:"error,omitempty"`
}

func decodeValidateCreateUserRequest(appConfig *conf.AppConfig) func(context.Context, *http.Request) (interface{}, error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var body createUserRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}

		// if err := appConfig.Validator.Struct(body); err != nil {
		// 	return nil, err
		// }

		return body, nil
	}
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case cargo.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func makeCreateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)

		user, err := s.RegisterUser(req.FullName, req.Username, req.Email, req.Password)
		if err != nil {
			return nil, err
		}

		return createUserResponse{
			user.ID,
			user.FullName,
			user.Username,
			user.Email,
			user.CreatedAt,
			user.UpdatedAt,
			nil,
		}, nil
	}
}
