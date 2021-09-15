package client

import (
	"context"
	"fmt"
	"github.com/ncostamagna/axul_user/pkg/grpc/userpb"
	"google.golang.org/grpc"
)

// Transport object
type Transport interface {
	GetAuth(id, token string) (int32, error)
}

type clientGRPC struct {
	client userpb.AuthServiceClient
}

// ClientType properties
type ClientType int

const (
	// HTTP transport type
	HTTP ClientType = iota

	// Socket transport type
	Socket

	// GRPC transport type
	GRPC
)

func NewClient(baseURL, token string, ct ClientType) Transport {

	switch ct {
	case GRPC:
		opts := grpc.WithInsecure()
		cc, err := grpc.Dial(baseURL, opts)
		if err != nil {
			panic(fmt.Sprintf("could not connect: %v", err))
		}

		return &clientGRPC{
			client: userpb.NewAuthServiceClient(cc),
		}
	}

	panic("Protocol hasn't been implement")
	return nil

}

func (c *clientGRPC) GetAuth(id, token string) (int32, error) {
	authReq := &userpb.AuthReq{
		Id:    id,
		Token: token,
	}

	ctx := context.Background()
	req, err := c.client.GetAuth(ctx, authReq)

	if err != nil {
		return 0, err
	}

	return req.Authorization, nil
}
