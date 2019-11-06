package transport

import (
	"context"
	"time"

	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/service"
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
	res.ResponsedAt = time.Now().Format(comtype.TimeLayout)

	res.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.writer.WriteHeader(res.StatusCode)
	json.NewEncoder(res.writer).Encode(res)
}

func (res *response) success(statusCode int, data interface{}) {
	res.Data = data
	res.response(statusCode)
}

func (res *response) fail(statusCode int, err *comtype.CommonError) {
	if len(err.Content) == 0 {
		res.Error = err.Error()
	} else {
		res.Error = err.Content
	}

	res.response(statusCode)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	result := &response{writer: w}

	switch err.(type) {
	case *comtype.CommonError:
		comErr := err.(*comtype.CommonError)

		switch comErr.Code {
		case comtype.ErrDataValidationFail:
			result.fail(http.StatusBadRequest, comErr)
			return
		case comtype.ErrDataNotFound:
			result.fail(http.StatusNotFound, comErr)
			return
		case comtype.ErrInvalidCredential:
			result.fail(http.StatusForbidden, comErr)
			return
		default:
			result.fail(http.StatusBadRequest, comErr)
			return
		}

	default:
		result.response(http.StatusInternalServerError)
		return
	}
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := &response{writer: w}
	result.success(http.StatusOK, data)

	return nil
}

type logErrorHandler struct {
	logger *log.Logger
}

func newLogErrorHandler(logger *log.Logger) *logErrorHandler {
	return &logErrorHandler{
		logger: logger,
	}
}

func (h *logErrorHandler) Handle(_ context.Context, err error) {
	switch err.(type) {
	case *comtype.CommonError:
		comErr := err.(*comtype.CommonError)
		if comErr.Code > 100 {
			h.logger.Error(comErr.Info())
		}

	default:
		h.logger.Error(err)
		return
	}
}

func decorateEndpointContext(s *service.AppService, c *conf.AppConfig) kithttp.RequestFunc {
	return func(ctx context.Context, _ *http.Request) context.Context {
		ctx = context.WithValue(ctx, comtype.CommonKeyRequestContext, s)
		ctx = context.WithValue(ctx, comtype.CommonKeyAppConfiguration, c)

		return ctx
	}
}
