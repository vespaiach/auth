package ep

import (
	"context"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/keymgr"
	"time"
)

type AddingKey struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
}

type Key struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Keys struct {
	Records []*Key `json:"records"`
	Total   int64  `json:"total"`
	Page    int64  `json:"page"`
	PerPage int64  `json:"per_page"`
}

type QueryingKey struct {
	Name    string
	Sort    string
	Page    int64
	PerPage int64
}

type ModifyingKey struct {
	Lookup string
	Key    string `json:"key"`
	Desc   string `json:"desc"`
}

type AddingKeyToBunch struct {
	Key   string
	Bunch string `json:"bunch"`
}

func AddingKeyEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	keych := make(chan *keymgr.Key)
	keyserv := ctx.Value(common.KeyManagementService).(keymgr.Service)

	go func() {
		req, ok := request.(*AddingKey)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		id, err := keyserv.AddKey(req.Key, req.Desc)
		if err != nil {
			erch <- err
			return
		}

		key, err := keyserv.GetKey(id)
		if err != nil {
			erch <- err
			return
		}
		keych <- key
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case key := <-keych:
		return &Key{
			key.ID,
			key.Key,
			key.Desc,
			key.CreatedAt,
			key.UpdatedAt,
		}, nil
	}
}

func ModifyingKeyEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	keych := make(chan *keymgr.Key)
	keyserv := ctx.Value(common.KeyManagementService).(keymgr.Service)

	go func() {
		req, ok := request.(*ModifyingKey)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		key, err := keyserv.GetKeyByName(req.Lookup)
		if err != nil {
			erch <- err
			return
		}
		if key == nil {
			erch <- common.ErrKeyNotFound
			return
		}

		err = keyserv.ModifyKey(key.ID, req.Key, req.Desc)
		if err != nil {
			erch <- err
			return
		}

		key, err = keyserv.GetKey(key.ID)
		if err != nil {
			erch <- err
			return
		}
		keych <- key
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case key := <-keych:
		return &Key{
			key.ID,
			key.Key,
			key.Desc,
			key.CreatedAt,
			key.UpdatedAt,
		}, nil
	}
}

func GettingKeyEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	keych := make(chan *keymgr.Key)
	keyserv := ctx.Value(common.KeyManagementService).(keymgr.Service)

	go func() {
		name, ok := request.(string)
		if !ok || len(name) == 0 {
			erch <- common.ErrWrongInputDatatype
			return
		}

		key, err := keyserv.GetKeyByName(name)
		if err != nil {
			erch <- err
			return
		}
		keych <- key
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case key := <-keych:
		return &Key{
			key.ID,
			key.Key,
			key.Desc,
			key.CreatedAt,
			key.UpdatedAt,
		}, nil
	}
}

func QueryingKeyEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	keych := make(chan []*keymgr.Key)
	keyserv := ctx.Value(common.KeyManagementService).(keymgr.Service)
	params, ok := request.(*QueryingKey)

	var total int64

	go func() {
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		if params.PerPage == 0 {
			params.PerPage = common.Take
		}

		if params.Page == 0 {
			params.Page = 1
		}

		records, count, err := keyserv.QueryKeys(params.Page, params.PerPage, params.Name, params.Sort)
		if err != nil {
			erch <- err
			return
		}
		total = count
		keych <- records
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case lst := <-keych:
		rows := make([]*Key, 0, len(lst))
		for _, row := range lst {
			rows = append(rows, &Key{
				row.ID,
				row.Key,
				row.Desc,
				row.CreatedAt,
				row.UpdatedAt,
			})
		}
		return &Keys{
			rows,
			total,
			params.Page,
			params.PerPage,
		}, nil
	}
}

func AddingKeyToBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	idch := make(chan int64)
	keyserv := ctx.Value(common.KeyManagementService).(keymgr.Service)

	go func() {
		req, ok := request.(*AddingKeyToBunch)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		bunchKeyID, err := keyserv.AddKeyToBunch(req.Key, req.Bunch)
		if err != nil {
			erch <- err
			return
		}

		idch <- bunchKeyID
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case id := <-idch:
		return id, nil
	}
}
