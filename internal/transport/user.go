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

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakeUsersHandlers publish endpoints through http
func MakeUsersHandlers(r *mux.Router, s *service.AppService, c *conf.AppConfig, log *log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorHandler(newLogErrorHandler(log)),
		kithttp.ServerBefore(decorateEndpointContext(s, c)),
	}

	registerUserHandler := kithttp.NewServer(
		endpoint.MakeRegisterUserEndpoint,
		decodeRegisterUserRequest,
		encodeResponse,
		opts...,
	)

	verifyLoginHandler := kithttp.NewServer(
		endpoint.MakeVerifyUserEndpoint,
		decodeVerifyUserRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/v1/users/login", verifyLoginHandler).Methods("POST")
	r.Handle("/v1/users/register", registerUserHandler).Methods("POST")

	return r
}

func decodeRegisterUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("decodeRegisterUserRequest:", err)
		return nil, err
	}

	return body, nil
}

func decodeVerifyUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.VerifyUserRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("decodeVerifyUserRequest:", err)
		return nil, err
	}
	body.RemoteAddr = r.RemoteAddr
	body.XForwardedFor = r.Header.Get("X-Forwarded-For")
	body.XRealIP = r.Header.Get("X-Real-Ip")
	body.UserAgent = r.Header.Get("User-Agent")

	return body, nil
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var body endpoint.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("decodeUpdateUserRequest:", err)
		return nil, err
	}

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return nil, comtype.NewCommonError(err, "decodeUpdateUserRequest", comtype.ErrBadRequest, nil)
	}
	body.ID = id

	return body, nil
}

// func decodeChangeUserPasswordRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
// 	var body endpoint.ChangeUserPasswordRequest

// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		log.Error("decodeUpdateUserRequest:", err)
// 		return nil, err
// 	}

// 	params := mux.Vars(r)
// 	id, err := strconv.ParseInt(params["id"], 10, 64)
// 	if err != nil {
// 		return nil, comtype.NewCommonError(err, "decodeUpdateUserRequest", comtype.ErrBadRequest, nil)
// 	}
// 	body.ID = id

// 	return body, nil
// }
