package tp

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/auth/jwt"
	kith "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/ep"
	"github.com/vespaiach/auth/pkg/keymgr"
	"net/http"
	"strconv"
)

// MakeKeyHandlers creates key endpoints through http
func MakeKeyHandlers(r *mux.Router, appConfig *cf.AppConfig, serv keymgr.Service) http.Handler {
	opts := []kith.ServerOption{
		kith.ServerErrorEncoder(encodeError),
		kith.ServerBefore(jwt.HTTPToContext()),
		kith.ServerBefore(addToContext(appConfig, common.AppConfigContextKey)),
		kith.ServerBefore(addToContext(serv, common.KeyManagementService)),
	}

	addKeyHandler := kith.NewServer(
		ep.AddingKeyEndPoint,
		decodeAddingKeyRequest,
		encodeResponse,
		opts...,
	)

	getKeyHandler := kith.NewServer(
		ep.GettingKeyEndPoint,
		decodeGettingKeyRequest,
		encodeResponse,
		opts...,
	)

	modifyKeyHandler := kith.NewServer(
		ep.ModifyingKeyEndPoint,
		decodeModifyingKeyRequest,
		encodeResponse,
		opts...,
	)

	queryKeyHandler := kith.NewServer(
		ep.QueryingKeyEndPoint,
		decodeQueryingKeyRequest,
		encodeResponse,
		opts...,
	)

	addKeyToBunchHandler := kith.NewServer(
		ep.AddingKeyToBunchEndPoint,
		decodeAddingKeyToBunchRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/keys", addKeyHandler).Methods("POST")
	r.Handle("/v1/keys/{key}", getKeyHandler).Methods("GET")
	r.Handle("/v1/keys/{key}", modifyKeyHandler).Methods("PATCH")
	r.Handle("/v1/keys", queryKeyHandler).Methods("GET")
	r.Handle("/v1/keys/{key}/bunch", addKeyToBunchHandler).Methods("POST")

	return r
}

func decodeAddingKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.AddingKey)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func decodeGettingKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["key"], nil
}

func decodeModifyingKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)

	data := new(ep.ModifyingKey)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	data.Lookup = params["key"]

	return data, nil
}

func decodeQueryingKeyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := r.URL.Query()
	data := &ep.QueryingKey{}

	name, nok := params["name"]
	if nok && len(name) > 0 {
		data.Name = name[0]
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

func decodeAddingKeyToBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)

	data := new(ep.AddingKeyToBunch)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	data.Key = params["key"]

	return data, nil
}
