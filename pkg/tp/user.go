package tp

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/ep"
	"github.com/vespaiach/auth/pkg/listing"
)

// MakeUserHandlers creates user endpoints through tp
func MakeUserHandlers(r *mux.Router, appConfig *cf.AppConfig, listingServ listing.Service, addServ adding.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(addServiceToContext(appConfig, common.AppConfigContextKey)),
		kithttp.ServerBefore(addServiceToContext(listingServ, common.ListingServiceContextKey)),
		kithttp.ServerBefore(addServiceToContext(addServ, common.AddingServiceContextKey)),
		kithttp.ServerBefore(jwt.HTTPToContext()),
	}

	getUserHandler := kithttp.NewServer(
		ep.GetUser,
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)

	addUserHandler := kithttp.NewServer(
		ep.AddUser,
		decodeAddUserRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/users/{id:[0-9]+}", getUserHandler).Methods("GET")
	r.Handle("/v1/users", addUserHandler).Methods("POST")

	return r
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func decodeAddUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.AddingUser)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
