package contacts

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

//NewHTTPServer is a server handler
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/contacts/", httptransport.NewServer(
		endpoints.Create,
		decodeGetReq,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/contacts/", httptransport.NewServer(
		endpoints.Update,
		decodeGetReq,
		encodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/contacts/", httptransport.NewServer(
		endpoints.GetAll,
		decodeGetReq,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/contacts/{id}", httptransport.NewServer(
		endpoints.Get,
		decodeGetReq,
		encodeResponse,
		opts...,
	)).Methods("GET")

	return r

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
