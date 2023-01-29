package contact

import (
	"context"
	"fmt"
	"github.com/digitalhouse-tech/go-lib-kit/meta"
	"github.com/digitalhouse-tech/go-lib-kit/response"
	"github.com/go-kit/kit/endpoint"
	"strconv"
	"time"
)

const (
	layoutISO = "2006-01-02 15:04:05"
)

// Endpoints struct
type (
	StoreReq struct {
		Auth      Authentication
		ID        uint   `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`
		Birthday  string `json:"birthday"`
	}

	GetReq struct {
		Auth Authentication
		ID   string
	}

	GetAllReq struct {
		Auth     Authentication
		Days     int64
		Birthday string
		Name     string
		Month    int16
		Limit    int
		Page     int
	}

	Authentication struct {
		ID    string
		Token string
	}

	Endpoints struct {
		Create endpoint.Endpoint
		Update endpoint.Endpoint
		Get    endpoint.Endpoint
		GetAll endpoint.Endpoint
		Alert  endpoint.Endpoint
	}
)

// MakeEndpoints handler endpoints
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Alert:  makeAlertEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreReq)
		if err := s.authorization(ctx, req.Auth.ID, req.Auth.Token); err != nil {
			return nil, response.Unauthorized(err.Error())
		}

		birthday, err := time.Parse(layoutISO, fmt.Sprintf("%s 17:00:00", req.Birthday))

		if err != nil {
			return nil, err
		}

		c, err := s.Create(ctx, req.Firstname, req.Lastname, req.Nickname, req.Gender, req.Phone, birthday)

		if err != nil {
			return nil, err
		}

		return response.Created("success", c, nil, nil), nil

	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)
		fmt.Println(req.Auth.ID, req.Auth.Token)
		if err := s.authorization(ctx, req.Auth.ID, req.Auth.Token); err != nil {
			return nil, response.Unauthorized(err.Error())
		}

		f := Filter{
			Name:      req.Name,
			Month:     req.Month,
			firstDate: time.Now().UTC(),
		}

		if req.Birthday != "" {
			days, err := strconv.Atoi(req.Birthday)
			if err != nil {
				return nil, response.BadRequest("Invalid birthday format in Query String")
			}

			f.Birthday = &days
		}

		if req.Days > 0 {
			f.RangeDays = &req.Days
		}

		count, err := s.Count(ctx, f)
		fmt.Println(count)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta := meta.New(req.Page, req.Limit, count)

		cs, err := s.GetAll(ctx, f, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success", cs, meta, nil), nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		if err := s.authorization(ctx, req.Auth.ID, req.Auth.Token); err != nil {
			return nil, response.Unauthorized(err.Error())
		}

		contact, err := s.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return response.OK("Success", contact, nil, nil), nil
	}
}

func makeAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)
		if err := s.authorization(ctx, req.Auth.ID, req.Auth.Token); err != nil {
			return nil, response.Unauthorized(err.Error())
		}

		fmt.Println(req)
		cs, err := s.Alert(ctx, req.Birthday)
		if err != nil {
			return nil, err
		}

		return response.OK("Success", cs, nil, nil), nil
	}
}
