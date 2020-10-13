package emails

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

	r.Handle("/emails/", httptransport.NewServer(
		endpoints.GetAll,
		decodeGetReq,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/emails/send", httptransport.NewServer(
		endpoints.Send,
		decodeEmailReq,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/emails/{id}/resend", httptransport.NewServer(
		endpoints.Resend,
		decodeResendReq,
		encodeResponse,
		opts...,
	)).Methods("POST")

	return r

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
