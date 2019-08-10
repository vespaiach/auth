package endpoint

import (
	"net/http"

	"github.com/vespaiach/auth/internal/conf"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/internal/service"
)

// MakeHandler returns a handler for the booking service.
func MakeHandler(us service.UserService, appConfig *conf.AppConfig) http.Handler {
	opts := []kithttp.ServerOption{}

	createUserHandler := kithttp.NewServer(
		makeCreateUserEndpoint(us),
		decodeValidateCreateUserRequest(appConfig),
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/user/v1/create", createUserHandler).Methods("POST")

	return r
}
