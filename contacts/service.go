package contacts

import (
	"context"
)

//Service interface
type Service interface {
	Create(ctx context.Context) (string, error)
	Update(ctx context.Context) (*Contact, error)
	Get(ctx context.Context) (Contact, error)
	GetAll(ctx context.Context) ([]Contact, error)
}
