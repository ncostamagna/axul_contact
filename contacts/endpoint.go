package contacts

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//Endpoints struct
type Endpoints struct {
	Create endpoint.Endpoint
	Update endpoint.Endpoint
	Get    endpoint.Endpoint
	GetAll endpoint.Endpoint
}

//MakeEndpoints handler endpoints
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}
