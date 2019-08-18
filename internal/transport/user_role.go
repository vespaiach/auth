package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/endpoint"
	"github.com/vespaiach/auth/internal/service"

	ep "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakeUserRoleHandlers creates role endpoints through http
func MakeUserRoleHandlers(r *mux.Router, s *service.AppService, c *conf.AppConfig, log *log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorHandler(newLogErrorHandler(log)),
		kithttp.ServerBefore(decorateEndpointContext(s, c)),
	}

	createUserRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("create_user_role"), jwtParser)(endpoint.MakeCreateUserRoleEndpoint),
		decodeCreateUserRoleRequest,
		encodeResponse,
		opts...,
	)

	deleteUserRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("delete_user_role"), jwtParser)(endpoint.MakeDeleteUserRoleEndpoint),
		decodeDeleteUserRoleRequest,
		encodeResponse,
		opts...,
	)

	getUserRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("get_user_role"), jwtParser)(endpoint.MakeGetUserRoleEndpoint),
		decodeGetUserRoleRequest,
		encodeResponse,
		opts...,
	)

	queryUserRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("query_user_role"), jwtParser)(endpoint.MakeQueryUserRoleEndpoint),
		decodeQueryUserRoleRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/user_roles", createUserRoleHandler).Methods("POST")
	r.Handle("/v1/user_roles/delete", deleteUserRoleHandler).Methods("POST")
	r.Handle("/v1/user_roles/{id:[0-9]+}", getUserRoleHandler).Methods("GET")
	r.Handle("/v1/user_roles", queryUserRoleHandler).Methods("GET")

	return r
}

func decodeCreateUserRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.CreateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeCreateUserRoleRequest", comtype.ErrBadRequest, nil)
	}

	return body, nil
}

func decodeDeleteUserRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.DeleteUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeDeleteUserRoleRequest", comtype.ErrBadRequest, nil)
	}

	return body, nil
}

func decodeGetUserRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeGetUserRoleRequest", comtype.ErrBadRequest, nil)
	}

	return id, nil
}

func decodeQueryUserRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var userID, roleID int64

	take, err := strconv.Atoi(r.FormValue("take"))
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeQueryUserRoleRequest", comtype.ErrBadRequest, nil)
	}

	if len(r.FormValue("user_id")) > 0 {
		userID, err = strconv.ParseInt(r.FormValue("user_id"), 10, 64)
		if err != nil {
			return nil, comtype.NewCommonError(err, "decodeQueryUserRoleRequest", comtype.ErrBadRequest, nil)
		}
	}

	if len(r.FormValue("role_id")) > 0 {
		roleID, err = strconv.ParseInt(r.FormValue("role_id"), 10, 64)
		if err != nil {
			return nil, comtype.NewCommonError(err, "decodeQueryUserRoleRequest", comtype.ErrBadRequest, nil)
		}
	}

	return endpoint.QueryUserRoleRequest{
		take,
		userID,
		roleID,
	}, nil
}
