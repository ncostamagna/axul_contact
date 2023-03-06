package handler

import (
	"context"
	"encoding/json"
	"github.com/digitalhouse-tech/go-lib-kit/response"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ncostamagna/axul_contact/internal/contact"
	"net/http"
	"strconv"
)

// NewHTTPServer is a server handler
func NewHTTPServer(ctx context.Context, endpoints contact.Endpoints) http.Handler {
	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.POST("/contacts", gin.WrapH(httptransport.NewServer(
		endpoints.Create,
		decodeCreateContact,
		encodeResponse,
		opts...,
	)))

	r.GET("/contacts", gin.WrapH(httptransport.NewServer(
		endpoints.GetAll,
		decodeGetAll,
		encodeResponse,
		opts...,
	)))

	r.GET("/contacts/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoints.Get,
		decodeGetContact,
		encodeResponse,
		opts...,
	)))

	r.PUT("/contacts/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoints.Update,
		nil,
		encodeResponse,
		opts...,
	)))

	r.DELETE("/contacts/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		nil,
		decodeCreateContact,
		encodeResponse,
		opts...,
	)))

	r.POST("/contacts/alert", gin.WrapH(httptransport.NewServer(
		endpoints.Alert,
		decodeGetAll,
		encodeResponse,
		opts...,
	)))

	return r

}

func ginDecode(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "params", c.Params)
	c.Request = c.Request.WithContext(ctx)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateContact(_ context.Context, r *http.Request) (interface{}, error) {
	var req contact.StoreReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	v := r.URL.Query()
	req.Auth.ID = v.Get("userid")
	req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func decodeGetContact(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := ctx.Value("params").(gin.Params)
	req := contact.GetReq{
		ID: pp.ByName("id"),
	}

	qs := r.URL.Query()
	req.Auth.ID = qs.Get("userid")
	req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func decodeGetAll(_ context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	d, _ := strconv.ParseInt(v.Get("days"), 0, 64)
	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	m, _ := strconv.ParseInt(v.Get("month"), 0, 64)
	req := contact.GetAllReq{
		Birthday: v.Get("birthday"),
		Days:     d,
		Month:    int16(m),
		Name:     v.Get("name"),
		Limit: limit,
		Page: page,
	}

	req.Auth.ID = v.Get("userid")
	req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
