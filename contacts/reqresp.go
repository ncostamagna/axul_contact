package contacts

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ncostamagna/response"
)

type (
	ContactRequest struct {
		ID        uint   `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`
		Birthday  string `json:"birthday"`
	}

	getRequest struct {
		id       string
		birthday string
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(r.GetStatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateContact(ctx context.Context, r *http.Request) (interface{}, error) {

	var req ContactRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetContact(ctx context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	req := getRequest{
		birthday: v.Get("birthday"),
	}
	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := response.NewResponse(err.Error(), 500, "", nil)
	w.WriteHeader(resp.GetStatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
