package contact

/*
import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "github.com/ncostamagna/axul_contact/templatespb"

	"google.golang.org/grpc"
)

type Transport interface {
	GetTemplate(id uint)
}

// ClientType -
type ClientType int

const (
	// HTTP transport type
	HTTP ClientType = iota

	// Grpc transport type
	GRPC
)

func NewClient(url, key string, ct ClientType) Transport {

	switch ct {
	case GRPC:
		fmt.Println("Entra")
		opts := grpc.WithInsecure()
		cc, err := grpc.Dial(url, opts)
		if err != nil {
			panic(fmt.Sprintf("could not connect: %v", err))
		}

		return &clientGRPC{
			client: pb.NewTemplatesServiceClient(cc),
		}
	}

	panic("Protocol hasn't been implement")

}

type clientGRPC struct {
	client pb.TemplatesServiceClient
}

func (f *clientGRPC) GetTemplate(id uint) {
	tempRequest := &pb.TemplateRequest{
		Id: uint32(id),
	}

	res, err := f.client.GetTemplate(context.Background(), tempRequest)

	if err != nil {
		log.Fatalf("error RPC: %v", err)
	}
	log.Printf("Response: %v", res)

}

func GetTemplateHTTP(id uint) {
	_, _ = http.Get(fmt.Sprintf("http://localhost:4000/templates/%d", id))


}
*/
