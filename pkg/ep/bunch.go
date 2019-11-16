package ep

import (
	"context"
	"database/sql"
	"github.com/vespaiach/auth/pkg/bunchmgr"
	"github.com/vespaiach/auth/pkg/common"
	"time"
)

type Bunch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddingBunch struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type ModifyingBunch struct {
	Lookup string
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Active *bool  `json:"active"`
}

type QueryingBunch struct {
	Name    string
	Active  sql.NullBool
	Sort    string
	Page    int64
	PerPage int64
}

type Bunches struct {
	Records []*Bunch `json:"records"`
	Total   int64    `json:"total"`
	Page    int64    `json:"page"`
	PerPage int64    `json:"per_page"`
}

type AddingKeysToBunch struct {
	Keys  []string `json:"keys"`
	Bunch string
}

func AddingBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	bch := make(chan *bunchmgr.Bunch)
	bserv := ctx.Value(common.BunchManagementService).(bunchmgr.Service)

	go func() {
		req, ok := request.(*AddingBunch)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		id, err := bserv.AddBunch(req.Name, req.Desc)
		if err != nil {
			erch <- err
			return
		}

		bunch, err := bserv.GetBunch(id)
		if err != nil {
			erch <- err
			return
		}
		bch <- bunch
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case b := <-bch:
		return &Bunch{
			b.ID,
			b.Name,
			b.Desc,
			b.Active.Bool,
			b.CreatedAt,
			b.UpdatedAt,
		}, nil
	}
}

func ModifyingBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	success := make(chan bool)
	bserv := ctx.Value(common.BunchManagementService).(bunchmgr.Service)

	go func() {
		req, ok := request.(*ModifyingBunch)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		bunch, err := bserv.GetBunchByName(req.Lookup)
		if err != nil {
			erch <- err
			return
		}
		if bunch == nil {
			erch <- common.ErrBunchNotFound
			return
		}

		active := sql.NullBool{}
		if req.Active != nil {
			active.Bool = *req.Active
			active.Valid = true
		}
		err = bserv.ModifyBunch(bunch.ID, req.Name, req.Desc, active)
		if err != nil {
			erch <- err
			return
		}
		success <- true
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case <-success:
		return true, nil
	}
}

func GettingBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	bch := make(chan *bunchmgr.Bunch)
	bserv := ctx.Value(common.BunchManagementService).(bunchmgr.Service)

	go func() {
		name, ok := request.(string)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		bunch, err := bserv.GetBunchByName(name)
		if err != nil {
			erch <- err
			return
		}
		if bunch == nil {
			erch <- common.ErrBunchNotFound
			return
		}

		bch <- bunch
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case b := <-bch:
		return &Bunch{
			b.ID,
			b.Name,
			b.Desc,
			b.Active.Bool,
			b.CreatedAt,
			b.UpdatedAt,
		}, nil
	}
}

func QueryingBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	bch := make(chan []*bunchmgr.Bunch)
	bserv := ctx.Value(common.BunchManagementService).(bunchmgr.Service)
	params, ok := request.(*QueryingBunch)

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

		records, count, err := bserv.QueryBunches(params.Page, params.PerPage, params.Name, params.Active, params.Sort)
		if err != nil {
			erch <- err
			return
		}
		total = count
		bch <- records
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case lst := <-bch:
		rows := make([]*Bunch, 0, len(lst))
		for _, row := range lst {
			rows = append(rows, &Bunch{
				row.ID,
				row.Name,
				row.Desc,
				row.Active.Bool,
				row.CreatedAt,
				row.UpdatedAt,
			})
		}
		return &Bunches{
			rows,
			total,
			params.Page,
			params.PerPage,
		}, nil
	}
}

func AddingKeysToBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	qch := make(chan bool)
	bserv := ctx.Value(common.BunchManagementService).(bunchmgr.Service)

	go func() {
		req, ok := request.(*AddingKeysToBunch)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		err := bserv.AddKeysToBunch(req.Bunch, req.Keys)
		if err != nil {
			erch <- err
			return
		}

		qch <- true
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case <-qch:
		return true, nil
	}
}

func GettingKeysInBunchEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	kch := make(chan []*bunchmgr.Key)
	bserv := ctx.Value(common.BunchManagementService).(bunchmgr.Service)

	go func() {
		name, ok := request.(string)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		keys, err := bserv.GetKeyInBunch(name)
		if err != nil {
			erch <- err
			return
		}

		kch <- keys
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case lst := <-kch:
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
		return rows, nil
	}
}
