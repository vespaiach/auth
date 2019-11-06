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

	"github.com/go-kit/kit/auth/jwt"
	ep "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakeActionHandlers creates action endpoints through http
func MakeActionHandlers(r *mux.Router, s *service.AppService, c *conf.AppConfig, log *log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorHandler(newLogErrorHandler(log)),
		kithttp.ServerBefore(decorateEndpointContext(s, c)),
		kithttp.ServerBefore(jwt.HTTPToContext()),
	}

	createActionHandler := kithttp.NewServer(
		ep.Chain(jwtParser, checkAuth("create_action"))(endpoint.MakeCreateActionEndpoint),
		decodeCreateActionRequest,
		encodeResponse,
		opts...,
	)

	updateActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("update_action"), jwtParser)(endpoint.MakeUpdateActionEndpoint),
		decodeUpdateActionRequest,
		encodeResponse,
		opts...,
	)

	getActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("get_action"), jwtParser)(endpoint.MakeGetActionEndpoint),
		decodeGetActionRequest,
		encodeResponse,
		opts...,
	)

	queryActionHandler := kithttp.NewServer(
		ep.Chain(checkAuth("query_action"), jwtParser)(endpoint.MakeQueryActionEndpoint),
		decodeQueryActionRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/actions", createActionHandler).Methods("POST")
	r.Handle("/v1/actions/{id:[0-9]+}", updateActionHandler).Methods("PATCH")
	r.Handle("/v1/actions/{id:[0-9]+}", getActionHandler).Methods("GET")
	r.Handle("/v1/actions", queryActionHandler).Methods("GET")

	return r
}

func decodeCreateActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.CreateActionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeCreateActionRequest", comtype.ErrBadRequest, nil)
	}

	return body, nil
}

func decodeUpdateActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)
	var body endpoint.UpdateActionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateActionRequest", comtype.ErrBadRequest, nil)
	}

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateActionRequest", comtype.ErrBadRequest, nil)
	}
	body.ID = id

	return body, nil
}

func decodeGetActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateActionRequest", comtype.ErrBadRequest, nil)
	}

	return id, nil
}

func decodeQueryActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	take, err := strconv.Atoi(r.FormValue("take"))
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeQueryActionRequest", comtype.ErrBadRequest, nil)
	}
	actionName := r.FormValue("action_name")
	sortBy := r.FormValue("sort_by")
	var active *bool
	if a := r.FormValue("active"); len(a) > 0 {
		parsed, _ := strconv.ParseBool(a)
		active = &parsed
	}

	return endpoint.QueryActionRequest{
		take,
		actionName,
		active,
		sortBy,
	}, nil
}
