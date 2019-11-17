package tp

import (
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	kith "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/bunchmgr"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/ep"
	"github.com/vespaiach/auth/pkg/keymgr"
	"github.com/vespaiach/auth/pkg/usrmgr"
)

type route struct {
	name          string
	path          string
	method        string
	endpoint      endpoint.Endpoint
	middleware    []endpoint.Middleware
	encoder       http.EncodeResponseFunc
	decoder       http.DecodeRequestFunc
	authorization bool
}

// List of routes
var routes = []*route{
	&route{
		name:          "login",
		path:          "/login",
		method:        "POST",
		endpoint:      ep.IssueTokenEndpoint,
		middleware:    []endpoint.Middleware{ep.VerifyingUserMiddleware},
		encoder:       encodeResponse,
		decoder:       decodeVerifyingUserUserRequest,
		authorization: false,
	},
	&route{
		name:          "add_user",
		path:          "/users",
		method:        "POST",
		endpoint:      ep.AddingUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeAddingUserRequest,
		authorization: true,
	},
	&route{
		name:          "modify_user",
		path:          "/users/{name}",
		method:        "PATCH",
		endpoint:      ep.ModifyingUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeModifyingUserRequest,
		authorization: true,
	},
	&route{
		name:          "get_user",
		path:          "/users/{name}",
		method:        "GET",
		endpoint:      ep.GettingUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeGettingUserRequest,
		authorization: true,
	},
	&route{
		name:          "query_user",
		path:          "/users",
		method:        "GET",
		endpoint:      ep.QueryingUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeQueryingUserRequest,
		authorization: true,
	},
	&route{
		name:          "query_user",
		path:          "/users",
		method:        "GET",
		endpoint:      ep.QueryingUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeQueryingUserRequest,
		authorization: true,
	},
	&route{
		name:          "add_bunch_to_user",
		path:          "/users/{name}/bunches",
		method:        "POST",
		endpoint:      ep.AddingBunchesToUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeAddingBunchesToUserRequest,
		authorization: true,
	},
	&route{
		name:          "get_bunch_of_user",
		path:          "/users/{name}/bunches",
		method:        "GET",
		endpoint:      ep.GettingBunchesOfUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeGettingBunchesOfUserRequest,
		authorization: true,
	},
	&route{
		name:          "get_key_of_user",
		path:          "/users/{name}/keys",
		method:        "GET",
		endpoint:      ep.GettingKeysOfUserEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeGettingKeysOfUserRequest,
		authorization: true,
	},
	&route{
		name:          "add_bunch",
		path:          "/bunches",
		method:        "POST",
		endpoint:      ep.AddingBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeAddingBunchRequest,
		authorization: true,
	},
	&route{
		name:          "modify_bunch",
		path:          "/bunches/{name}",
		method:        "POST",
		endpoint:      ep.ModifyingBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeModifyingBunchRequest,
		authorization: true,
	},
	&route{
		name:          "get_bunch",
		path:          "/bunches/{name}",
		method:        "GET",
		endpoint:      ep.GettingBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeGettingBunchRequest,
		authorization: true,
	},
	&route{
		name:          "query_bunch",
		path:          "/bunches",
		method:        "GET",
		endpoint:      ep.QueryingBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeQueryingBunchRequest,
		authorization: true,
	},
	&route{
		name:          "add_keys_to_bunch",
		path:          "/bunches/{name}/keys",
		method:        "POST",
		endpoint:      ep.AddingKeysToBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeAddingKeysToBunchRequest,
		authorization: true,
	},
	&route{
		name:          "get_key_of_bunch",
		path:          "/bunches/{name}/keys",
		method:        "GET",
		endpoint:      ep.GettingKeysInBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeGettingKeysInBunchRequest,
		authorization: true,
	},
	&route{
		name:          "add_key",
		path:          "/keys",
		method:        "POST",
		endpoint:      ep.AddingKeyEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeAddingKeyRequest,
		authorization: true,
	},
	&route{
		name:          "modify_key",
		path:          "/keys/{key}",
		method:        "POST",
		endpoint:      ep.ModifyingKeyEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeModifyingKeyRequest,
		authorization: true,
	},
	&route{
		name:          "get_key",
		path:          "/keys/{key}",
		method:        "GET",
		endpoint:      ep.GettingKeyEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeGettingKeyRequest,
		authorization: true,
	},
	&route{
		name:          "query_key",
		path:          "/keys",
		method:        "GET",
		endpoint:      ep.QueryingKeyEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeQueryingKeyRequest,
		authorization: true,
	},
	&route{
		name:          "add_key_to_bunch",
		path:          "/keys/{key}/bunch",
		method:        "POST",
		endpoint:      ep.AddingKeyToBunchEndpoint,
		middleware:    nil,
		encoder:       encodeResponse,
		decoder:       decodeAddingKeyToBunchRequest,
		authorization: true,
	},
}

func makeHandler(r *route, opts []kith.ServerOption) *kith.Server {
	mids := make([]endpoint.Middleware, 0)

	if r.authorization {
		mids = append(mids, ep.TokenParserMiddleware, ep.KeyCheckerMiddleware(r.name))
	}

	if r.middleware != nil {
		mids = append(mids, r.middleware...)
	}

	if len(mids) > 0 {
		return kith.NewServer(
			endpoint.Chain(mids[0], mids[1:]...)(r.endpoint),
			r.decoder,
			r.encoder,
			opts...
		)
	}
	return kith.NewServer(
		r.endpoint,
		r.decoder,
		r.encoder,
		opts...
	)
}

func CreateRouter(appConfig *cf.AppConfig, userServ usrmgr.Service, bunchServ bunchmgr.Service, keyServ keymgr.Service) *mux.Router {
	router := mux.NewRouter()
	opts := []kith.ServerOption{
		kith.ServerErrorEncoder(encodeError),
		kith.ServerBefore(jwt.HTTPToContext()),
		kith.ServerBefore(addToContext(appConfig, common.AppConfigContextKey)),
		kith.ServerBefore(addToContext(userServ, common.UserManagementService)),
		kith.ServerBefore(addToContext(bunchServ, common.BunchManagementService)),
		kith.ServerBefore(addToContext(keyServ, common.KeyManagementService)),
	}

	for _, r := range routes {
		router.PathPrefix("/v1").
			Path(r.path).
			Methods(r.method).
			Handler(makeHandler(r, opts))
	}

	return router
}
