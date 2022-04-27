package contacts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ncostamagna/response"
)

type (
	ContactRequest struct {
		Auth      Authentication
		ID        uint   `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`
		Birthday  string `json:"birthday"`
	}

	getRequest struct {
		Auth Authentication
		id   string
	}

	getAllReq struct {
		Auth     Authentication
		days     int64
		birthday string
		name     string
		month    int16
	}

	Authentication struct {
		ID    string
		Token string
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(r.GetStatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateContact(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ContactRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	v := r.URL.Query()
	req.Auth.ID = v.Get("userid")
	req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func decodeGetContact(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := mux.Vars(r)

	req := getRequest{
		id: pp["id"],
	}

	qs := r.URL.Query()
	req.Auth.ID = qs.Get("userid")
	req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func decodeGetAll(ctx context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()
	fmt.Println(v)
	d, _ := strconv.ParseInt(v.Get("days"), 0, 64)
	fmt.Println("days: ", d)

	m, _ := strconv.ParseInt(v.Get("month"), 0, 64)
	req := getAllReq{
		birthday: v.Get("birthday"),
		days:     d,
		month:    int16(m),
		name:     v.Get("name"),
	}

	req.Auth.ID = v.Get("userid")
	req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := response.NewResponse(err.Error(), 500, "", nil)
	w.WriteHeader(resp.GetStatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
