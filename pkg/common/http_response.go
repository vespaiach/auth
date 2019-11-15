package common

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {
	Data        interface{} `json:"data,omitempty"`
	Error       interface{} `json:"error,omitempty"`
	Status      interface{} `json:"status"`
	StatusCode  int         `json:"code"`
	ResponsedAt string      `json:"Responsed_at"`
	writer      http.ResponseWriter
}

func (res *Response) Response(statusCode int) {
	res.StatusCode = statusCode
	res.Status = http.StatusText(res.StatusCode)
	res.ResponsedAt = time.Now().Format(TimeLayout)

	res.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.writer.WriteHeader(res.StatusCode)
	json.NewEncoder(res.writer).Encode(res)
}

func (res *Response) Success(statusCode int, data interface{}) {
	res.Data = data
	res.Response(statusCode)
}

func (res *Response) Fail(statusCode int, err error) {
	res.Error = err.Error()
	res.Response(statusCode)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	result := &Response{writer: w}

	result.Response(http.StatusInternalServerError)
	return
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result := &Response{writer: w}
	result.Success(http.StatusOK, data)

	return nil
}
