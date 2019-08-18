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

// MakeRoleHandlers creates role endpoints through http
func MakeRoleHandlers(r *mux.Router, s *service.AppService, c *conf.AppConfig, log *log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorHandler(newLogErrorHandler(log)),
		kithttp.ServerBefore(decorateEndpointContext(s, c)),
	}

	createRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("create_role"), jwtParser)(endpoint.MakeCreateRoleEndpoint),
		decodeCreateRoleRequest,
		encodeResponse,
		opts...,
	)

	updateRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("update_role"), jwtParser)(endpoint.MakeUpdateRoleEndpoint),
		decodeUpdateRoleRequest,
		encodeResponse,
		opts...,
	)

	getRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("get_role"), jwtParser)(endpoint.MakeGetRoleEndpoint),
		decodeGetRoleRequest,
		encodeResponse,
		opts...,
	)

	queryRoleHandler := kithttp.NewServer(
		ep.Chain(checkAuth("query_role"), jwtParser)(endpoint.MakeQueryRoleEndpoint),
		decodeQueryRoleRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/roles", createRoleHandler).Methods("POST")
	r.Handle("/v1/roles/{id:[0-9]+}", updateRoleHandler).Methods("PATCH")
	r.Handle("/v1/roles/{id:[0-9]+}", getRoleHandler).Methods("GET")
	r.Handle("/v1/roles", queryRoleHandler).Methods("GET")

	return r
}

func decodeCreateRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeCreateRoleRequest", comtype.ErrBadRequest, nil)
	}

	return body, nil
}

func decodeUpdateRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)
	var body endpoint.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateRoleRequest", comtype.ErrBadRequest, nil)
	}

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateRoleRequest", comtype.ErrBadRequest, nil)
	}
	body.ID = id

	return body, nil
}

func decodeGetRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateRoleRequest", comtype.ErrBadRequest, nil)
	}

	return id, nil
}

func decodeQueryRoleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	take, err := strconv.Atoi(r.FormValue("take"))
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeQueryRoleRequest", comtype.ErrBadRequest, nil)
	}
	roleName := r.FormValue("role_name")
	sortBy := r.FormValue("sort_by")
	var active *bool
	if a := r.FormValue("active"); len(a) > 0 {
		parsed, _ := strconv.ParseBool(a)
		active = &parsed
	}

	return endpoint.QueryRoleRequest{
		take,
		roleName,
		active,
		sortBy,
	}, nil
}
