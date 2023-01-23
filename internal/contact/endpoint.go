package contact

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalhouse-tech/go-lib-kit/response"
	"github.com/go-kit/kit/endpoint"
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
			return nil, err
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
			return nil, err
		}

		var contacts []Contact

		f := Filter{
			birthday: req.Birthday,
			days:     req.Days,
			name:     req.Name,
			month:    req.Month,
		}

		if err := s.GetAll(ctx, &contacts, f); err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success", contacts, nil, nil), nil
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
			return nil, response.InternalServerError(err.Error())
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
			return nil, err
		}

		var contacts []Contact
		fmt.Println(req)

		if err := s.Alert(ctx, &contacts, req.Birthday); err != nil {
			return nil, err
		}

		return response.OK("Success", contacts, nil, nil), nil
	}
}
