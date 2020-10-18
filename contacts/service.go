package contacts

import (
	"context"

	"github.com/ncostamagna/rerrors"
)

//Service interface
type Service interface {
	Create(ctx context.Context, contact *Contact) rerrors.RestErr
	Update(ctx context.Context) (*Contact, rerrors.RestErr)
	Get(ctx context.Context) (Contact, rerrors.RestErr)
	GetAll(ctx context.Context, contacts *[]Contact, birthday string) rerrors.RestErr
}
