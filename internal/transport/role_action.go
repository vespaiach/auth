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

// MakeRoleActionHandlers creates role endpoints through http
func MakeRoleActionHandlers(r *mux.Router, s *service.AppService, c *conf.AppConfig, log *log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorHandler(newLogErrorHandler(log)),
		kithttp.ServerBefore(decorateEndpointContext(s, c)),
	}

	createRoleActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("create_role_action"), jwtParser)(endpoint.MakeCreateRoleActionEndpoint),
		decodeCreateRoleActionRequest,
		encodeResponse,
		opts...,
	)

	deleteRoleActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("delete_role_action"), jwtParser)(endpoint.MakeDeleteRoleActionEndpoint),
		decodeDeleteRoleActionRequest,
		encodeResponse,
		opts...,
	)

	getRoleActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("get_role_action"), jwtParser)(endpoint.MakeGetRoleActionEndpoint),
		decodeGetRoleActionRequest,
		encodeResponse,
		opts...,
	)

	queryRoleActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("query_role_action"), jwtParser)(endpoint.MakeQueryRoleActionEndpoint),
		decodeQueryRoleActionRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/role_actions", createRoleActionHandler).Methods("POST")
	r.Handle("/v1/role_actions/delete", deleteRoleActionHandler).Methods("POST")
	r.Handle("/v1/role_actions/{id:[0-9]+}", getRoleActionHandler).Methods("GET")
	r.Handle("/v1/role_actions", queryRoleActionHandler).Methods("GET")

	return r
}

func decodeCreateRoleActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.CreateRoleActionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeCreateRoleActionRequest", comtype.ErrBadRequest, nil)
	}

	return body, nil
}

func decodeDeleteRoleActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.DeleteRoleActionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeDeleteRoleActionRequest", comtype.ErrBadRequest, nil)
	}

	return body, nil
}

func decodeGetRoleActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeGetRoleActionRequest", comtype.ErrBadRequest, nil)
	}

	return id, nil
}

func decodeQueryRoleActionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		roleID, actionID int64
		err              error
	)

	take, err := strconv.Atoi(r.FormValue("take"))
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeQueryRoleActionRequest", comtype.ErrBadRequest, nil)
	}

	if len(r.FormValue("action_id")) > 0 {
		actionID, err = strconv.ParseInt(r.FormValue("action_id"), 10, 64)
		if err != nil {
			return nil, comtype.NewCommonError(err, "decodeQueryRoleActionRequest", comtype.ErrBadRequest, nil)
		}
	}

	if len(r.FormValue("role_id")) > 0 {
		roleID, err = strconv.ParseInt(r.FormValue("role_id"), 10, 64)
		if err != nil {
			return nil, comtype.NewCommonError(err, "decodeQueryRoleActionRequest", comtype.ErrBadRequest, nil)
		}
	}

	return endpoint.QueryRoleActionRequest{
		take,
		roleID,
		actionID,
	}, nil
}
