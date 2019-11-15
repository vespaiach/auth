package keymgr

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

var errWrongData = errors.New("couldn't get data")

type addingKey struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
}

func makeAddingKeyEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addingKey)

		id, err := s.AddKey(req.Key, req.Desc)
		if err != nil {
			return nil, err
		}

		return id, nil
	}
}

type modifyingKey struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

func makeModifyKeyEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(modifyingKey)
		if !ok {
			return nil, errWrongData
		}

		success, err := s.ModifyKey(req.ID, req.Key)
		if err != nil {
			return nil, err
		}

		return success, nil
	}
}

func makeGetKeyEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, ok := request.(int64)
		if !ok {
			return nil, errWrongData
		}

		key, err := s.GetKey(id)
		if err != nil {
			return nil, err
		}

		return key, nil
	}
}
