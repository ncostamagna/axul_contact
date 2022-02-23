package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ncostamagna/axul_user/pkg/grpc/userpb"
	c "github.com/ncostamagna/streetflow/client"
	"google.golang.org/grpc"
)

// Transport object
type Transport interface {
	GetAuth(id, token string) (int32, error)
}

type clientGRPC struct {
	client userpb.AuthServiceClient
}

type clientHTTP struct {
	client c.RequestBuilder
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
	case HTTP:
		header := http.Header{}
		//header.Set("X-Api-Key", token)
		return &clientHTTP{
			client: c.RequestBuilder{
				Headers:        header,
				BaseURL:        baseURL,
				ConnectTimeout: 5000 * time.Millisecond,
				LogTime:        true,
			},
		}
	}

	panic("Protocol hasn't been implement")
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

func (c *clientHTTP) GetAuth(id, token string) (int32, error) {
	// le pega
	return 0, nil
}
