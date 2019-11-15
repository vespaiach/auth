package keymgr

import (
	"net/http"

	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandlers returns a handler for the key service.
func MakeHandlers(r *mux.Route) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	loadCargoHandler := kithttp.NewServer(
		makeLoadCargoEndpoint(bs),
		decodeLoadCargoRequest,
		encodeResponse,
		opts...,
	)

	listLocationsHandler := kithttp.NewServer(
		makeListLocationsEndpoint(bs),
		decodeListLocationsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/booking/v1/cargos", bookCargoHandler).Methods("POST")
	r.Handle("/booking/v1/cargos", listCargosHandler).Methods("GET")
	r.Handle("/booking/v1/cargos/{id}", loadCargoHandler).Methods("GET")
	r.Handle("/booking/v1/cargos/{id}/request_routes", requestRoutesHandler).Methods("GET")
	r.Handle("/booking/v1/cargos/{id}/assign_to_route", assignToRouteHandler).Methods("POST")
	r.Handle("/booking/v1/cargos/{id}/change_destination", changeDestinationHandler).Methods("POST")
	r.Handle("/booking/v1/locations", listLocationsHandler).Methods("GET")

	return r
}
