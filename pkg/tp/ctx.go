package tp

import (
	"context"
	"github.com/vespaiach/auth/pkg/common"
	"net/http"

	kith "github.com/go-kit/kit/transport/http"
)

import (
	"encoding/json"
	"time"
)

type response struct {
	Data        interface{} `json:"data,omitempty"`
	Error       interface{} `json:"error,omitempty"`
	Status      interface{} `json:"status"`
	StatusCode  int         `json:"code"`
	ResponsedAt string      `json:"responsed_at"`
	writer      http.ResponseWriter
}

func (res *response) response(statusCode int) {
	res.StatusCode = statusCode
	res.Status = http.StatusText(res.StatusCode)
	res.ResponsedAt = time.Now().Format(common.TimeLayout)

	res.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.writer.WriteHeader(res.StatusCode)
	json.NewEncoder(res.writer).Encode(res)
}

func (res *response) success(statusCode int, data interface{}) {
	res.Data = data
	res.response(statusCode)
}

func (res *response) fail(statusCode int, err error) {
	res.Error = err.Error()
	res.response(statusCode)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	result := &response{writer: w}
	switch err {
	case common.ErrWrongInputDatatype:
		result.fail(http.StatusInternalServerError, err)
		break
	case common.ErrWrongJWTToken, common.ErrMissingJWTToken, common.ErrWrongCredentials:
		result.fail(http.StatusUnauthorized, err)
		break
	case common.ErrNotAllowed:
		result.fail(http.StatusForbidden, err)
		break
	case common.ErrDuplicatedUsername, common.ErrEmailInvalid, common.ErrDuplicatedEmail,
		common.ErrUsernameInvalid, common.ErrDuplicatedBunch, common.ErrBunchNameInvalid,
		common.ErrKeyNameInvalid, common.ErrMissingHash, common.ErrDuplicatedKey,
		common.ErrPasswordMissing:
		result.fail(http.StatusBadRequest, err)
		break
	case common.ErrKeyNotFound, common.ErrUserNotFound, common.ErrBunchNotFound:
		result.fail(http.StatusNotFound, err)
		break
	default:
		result.fail(http.StatusInternalServerError, err)
		break
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := &response{writer: w}
	result.success(http.StatusOK, data)

	return nil
}

func addToContext(s interface{}, key common.ContextKey) kith.RequestFunc {
	return func(ctx context.Context, _ *http.Request) context.Context {
		return context.WithValue(ctx, key, s)
	}
}
