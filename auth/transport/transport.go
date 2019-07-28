package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/vespaiach/auth/service"
	"github.com/vespaiach/auth/store"

	"github.com/gorilla/mux"
)

// ErrInvalidPayload invalid payload
var ErrInvalidPayload = errors.New("invalid payload")

// MakeHandler returns a handler for the tracking service.
func MakeHandler(auth service.Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	loginHandler := kithttp.NewServer(
		makeVerifyLoginEndpoint(auth),
		decodeVerifyLoginRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/users/v1/login", loginHandler).Methods("POST")

	return r
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case store.ErrDataNotFound:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidPayload:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"err": err.Error(),
	})
}
