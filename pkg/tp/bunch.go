package tp

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/ep"
	"net/http"
	"strconv"
)

func decodeAddingBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.AddingBunch)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func decodeModifyingBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	data := new(ep.ModifyingBunch)
	params := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	data.Lookup = params["name"]

	return data, nil
}

func decodeGettingBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["name"], nil
}

func decodeQueryingBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := r.URL.Query()
	data := &ep.QueryingBunch{}

	name, nok := params["name"]
	if nok && len(name) > 0 {
		data.Name = name[0]
	}

	active, aok := params["active"]
	if aok && len(active) > 0 {
		b, err := strconv.ParseBool(active[0])
		if err != nil {
			return nil, err
		}
		data.Active = sql.NullBool{
			Bool:  b,
			Valid: true,
		}
	}

	sort, sok := params["sort"]
	if sok && len(sort) > 0 {
		data.Sort = sort[0]
	}

	page, pok := params["page"]
	if pok && len(page) > 0 {
		intPage, err := strconv.ParseInt(page[0], 10, 64)
		if err != nil {
			return nil, err
		}
		data.Page = intPage
	}

	perPage, ppok := params["per_page"]
	if ppok && len(perPage) > 0 {
		intPerPage, err := strconv.ParseInt(perPage[0], 10, 64)
		if err != nil {
			return nil, err
		}
		data.PerPage = intPerPage
	}

	return data, nil
}

func decodeAddingKeysToBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)

	data := new(ep.AddingKeysToBunch)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	data.Bunch = params["name"]

	return data, nil
}

func decodeGettingKeysInBunchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["name"], nil
}
