package tp

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/ep"
	"github.com/vespaiach/auth/pkg/listing"
	"net/http"
	"strconv"
)

// MakeUserHandlers creates user endpoints through tp
func MakeUserHandlers(r *mux.Router, listing listing.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(addServiceToContext(listing, common.ListingServiceContextKey)),
		kithttp.ServerBefore(jwt.HTTPToContext()),
	}

	getUserHandler := kithttp.NewServer(
		ep.GetUser,
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/users/{id:[0-9]+}", getUserHandler).Methods("GET")

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
