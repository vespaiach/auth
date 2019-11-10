package tp

import (
	"context"
	"encoding/json"
	kith "github.com/go-kit/kit/transport/http"
	"github.com/vespaiach/auth/pkg/common"
	"net/http"
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
	appErr, ok := err.(*common.AppErr)
	if ok && len(appErr.Payload) > 0 {
		res.Error = appErr.Payload
	} else {
		res.Error = err.Error()
	}
	res.response(statusCode)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	result := &response{writer: w}

	switch err.(type) {
	case *common.AppErr:
		appErr := err.(*common.AppErr)

		switch appErr.Code {
		case common.ErrGetData:
		case common.ErrExecData:
			result.fail(http.StatusInternalServerError, appErr)
			return
		case common.ErrDataNotFound:
			result.fail(http.StatusNotFound, appErr)
			return
		default:
			result.fail(http.StatusBadRequest, appErr)
			return
		}

	default:
		result.response(http.StatusInternalServerError)
		return
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := &response{writer: w}
	result.success(http.StatusOK, data)

	return nil
}

func addServiceToContext(s interface{}, key common.ContextKey) kith.RequestFunc {
	return func(ctx context.Context, _ *http.Request) context.Context {
		return context.WithValue(ctx, key, s)
	}
}
