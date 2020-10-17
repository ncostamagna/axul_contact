package contacts

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"gitlab.com/digitalhouse-dev/libraries-go/util.git/response"
)

//Endpoints struct
type Endpoints struct {
	Send   endpoint.Endpoint
	GetAll endpoint.Endpoint
	Resend endpoint.Endpoint
}

//MakeEndpoints handler endpoints
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Send:   makeSendEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Resend: makeResendEndpoint(s),
	}
}

// makeSendEndpoint endpoint
// @Summary Send Email
// @Tags email
// @Accept  json
// @Produce  json
// @Param email body SendMailRequest true "email to sender"
// @success 200 {object} sendMailResponseOK
// @Failure 400 {object} sendMailResponseError
// @Failure 500 {object} sendMailResponseError
// @Router /email/send [post]
func makeSendEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendMailRequest)

		ok, err := s.Send(ctx, req.Subject, req.Content.Type, req.Content.Value, req.From, req.ReplyTo, req.To, req.Cc, req.Bcc, req.SandboxMode, req.Attachments)

		if err != nil {
			resp := response.NewResponse(err.Message(), err.Status(), "", err.Message())
			return resp, nil
		}

		resp := response.NewResponse(ok, http.StatusOK, "", req)
		return resp, nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EmailRequest)

		emails, err := s.GetAll(ctx, req.WasSent)

		if err != nil {
			resp := response.NewResponse(err.Message(), err.Status(), "", err.Message())
			return resp, nil
		}

		resp := response.NewResponse("", http.StatusOK, "", emails)
		return resp, nil
	}
}

func makeResendEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EmailRequest)

		email, err := s.Resend(ctx, req.ID)

		if err != nil {
			resp := response.NewResponse(err.Message(), err.Status(), "", err.Message())
			return resp, nil
		}

		resp := response.NewResponse("", http.StatusOK, "", email)
		return resp, nil
	}
}
