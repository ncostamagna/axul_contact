package emails

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/digitalhouse-dev/libraries-go/util.git/response"
)

type (
	//SendMailRequest is a email sender request entity by POST
	SendMailRequest struct {
		To          []addressStruct `json:"to"`
		Cc          []addressStruct `json:"cc"`
		Bcc         []addressStruct `json:"bcc"`
		From        addressStruct   `json:"from"`
		ReplyTo     addressStruct   `json:"replyTo"`
		Subject     string          `json:"subject"`
		Content     contentStruct   `json:"content"`
		SandboxMode bool            `json:"sandboxMode"`
		Attachments []attachStruct  `json:"attachments"`
	}

	// EmailRequest is a email request entity with the following fields:
	//
	// `WasSent`: '0' -  Wasn't Sent, '1' - Was Sent.`ID`: int
	EmailRequest struct {
		WasSent string
		ID      uint
	}
)

type (
	addressStruct struct {
		Name    string `json:"name"`
		Address string `json:"email"`
	}

	contentStruct struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	attachStruct struct {
		Content     string `json:"content"`
		Type        string `json:"type"`
		FileName    string `json:"filename"`
		Disposition string `json:"disposition"`
		ContentID   string `json:"content_id"`
	}

	// ****************
	//   only swagger
	// ****************

	//nolint
	sendMailResponseOK struct {
		Message    string          `json:"message"`
		Code       string          `json:"code"`
		StatusCode int             `json:"-"`
		Data       SendMailRequest `json:"data"`
		Meta       struct {
			Limit       int `json:"limit"`
			CurrentPage int `json:"current"`
		} `json:"meta"`
	}

	//nolint
	sendMailResponseError struct {
		Message    string      `json:"message"`
		Code       string      `json:"code"`
		StatusCode int         `json:"-"`
		Errors     interface{} `json:"errors"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(r.GetStatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeEmailReq(ctx context.Context, r *http.Request) (interface{}, error) {

	var req SendMailRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetReq(ctx context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	req := EmailRequest{
		WasSent: v.Get("wasSent"),
	}
	return req, nil
}

func decodeResendReq(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	fmt.Println(vars["id"])
	fmt.Println(id)
	fmt.Println(uint(id))
	fmt.Println(err)

	req := EmailRequest{
		ID: uint(id),
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := response.NewResponse(err.Error(), 500, "", nil)
	w.WriteHeader(resp.GetStatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
