package tp

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kith "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/ep"
	"github.com/vespaiach/auth/pkg/usrmgr"
	"net/http"
	"strconv"
)

// MakeUserHandlers creates user endpoints through http
func MakeUserHandlers(r *mux.Router, appConfig *cf.AppConfig, serv usrmgr.Service) http.Handler {
	opts := []kith.ServerOption{
		kith.ServerErrorEncoder(encodeError),
		kith.ServerBefore(jwt.HTTPToContext()),
		kith.ServerBefore(addToContext(appConfig, common.AppConfigContextKey)),
		kith.ServerBefore(addToContext(serv, common.UserManagementService)),
	}

	addUserHandler := kith.NewServer(
		ep.AddingUserEndPoint,
		decodeAddingUserRequest,
		encodeResponse,
		opts...,
	)

	getUserHandler := kith.NewServer(
		ep.GettingUserEndPoint,
		decodeGettingUserRequest,
		encodeResponse,
		opts...,
	)

	modifyUserHandler := kith.NewServer(
		ep.ModifyingUserEndPoint,
		decodeModifyingUserRequest,
		encodeResponse,
		opts...,
	)

	queryUserHandler := kith.NewServer(
		ep.QueryingUserEndPoint,
		decodeQueryingUserRequest,
		encodeResponse,
		opts...,
	)

	addBunchesToUserHandler := kith.NewServer(
		ep.AddingBunchesToUserEndPoint,
		decodeAddingBunchesToUserRequest,
		encodeResponse,
		opts...,
	)

	getBunchesOfUserHandler := kith.NewServer(
		ep.GettingBunchesOfUserEndPoint,
		decodeGettingBunchesOfUserRequest,
		encodeResponse,
		opts...,
	)

	getKeysOfUserHandler := kith.NewServer(
		ep.GettingKeysOfUserEndPoint,
		decodeGettingKeysOfUserRequest,
		encodeResponse,
		opts...,
	)

	loginUserHandler := kith.NewServer(
		endpoint.Chain(ep.VerifyingUserMiddleware)(ep.IssueTokenEndPoint),
		decodeVerifyingUserUserRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/users", addUserHandler).Methods("POST")
	r.Handle("/v1/users/{name}", getUserHandler).Methods("GET")
	r.Handle("/v1/users/{name}", modifyUserHandler).Methods("PATCH")
	r.Handle("/v1/users", queryUserHandler).Methods("GET")
	r.Handle("/v1/users/{name}/bunches", addBunchesToUserHandler).Methods("POST")
	r.Handle("/v1/users/{name}/bunches", getBunchesOfUserHandler).Methods("GET")
	r.Handle("/v1/users/{name}/keys", getKeysOfUserHandler).Methods("GET")
	r.Handle("/v1/login", loginUserHandler).Methods("POST")

	return r
}

func decodeAddingUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.AddingUser)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func decodeGettingUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["name"], nil
}

func decodeModifyingUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.ModifyingUser)
	params := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	data.Lookup = params["name"]

	return data, nil
}

func decodeQueryingUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := r.URL.Query()
	data := &ep.QueryingUser{}

	name, nok := params["username"]
	if nok && len(name) > 0 {
		data.Username = name[0]
	}

	email, eok := params["email"]
	if eok && len(email) > 0 {
		data.Email = email[0]
	}

	active, aok := params["active"]
	if aok && len(active) > 0 {
		b, err := strconv.ParseBool(active[0])
		if err != nil {
			return nil, err
		}
		data.Active = sql.NullBool{
			Bool:  b,
			Valid: true,
		}
	}

	sort, sok := params["sort"]
	if sok && len(sort) > 0 {
		data.Sort = sort[0]
	}

	page, pok := params["page"]
	if pok && len(page) > 0 {
		intPage, err := strconv.ParseInt(page[0], 10, 64)
		if err != nil {
			return nil, err
		}
		data.Page = intPage
	}

	perPage, ppok := params["per_page"]
	if ppok && len(perPage) > 0 {
		intPerPage, err := strconv.ParseInt(perPage[0], 10, 64)
		if err != nil {
			return nil, err
		}
		data.PerPage = intPerPage
	}

	return data, nil
}

func decodeAddingBunchesToUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)

	data := new(ep.AddingBunchesToUser)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	data.Username = params["name"]

	return data, nil
}

func decodeGettingBunchesOfUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["name"], nil
}

func decodeGettingKeysOfUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["name"], nil
}

func decodeVerifyingUserUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.VerifyingUser)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
