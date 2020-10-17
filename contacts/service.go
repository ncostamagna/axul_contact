package contacts

import (
	"context"

	"gitlab.com/digitalhouse-dev/libraries-go/util.git/rerrors"
)

//Service interface
type Service interface {
	Create(ctx context.Context) (string, rerrors.RestErr)
	Update(ctx context.Context) (*Contact, rerrors.RestErr)
	Get(ctx context.Context) (Contact, rerrors.RestErr)
	GetAll(ctx context.Context) ([]Contact, rerrors.RestErr)
}
