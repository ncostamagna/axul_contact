package contacts

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type service struct {
	repo   Repository
	logger log.Logger
}

type updateCb func(uint, time.Time) error

//NewService is a service handler
func NewService(repo Repository, logger log.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

//Create service
func (s service) Create(ctx context.Context) (string, error) {

	return "Success", nil
}

func (s service) Update(ctx context.Context) (*Contact, error) {

	contact := Contact{}

	return &contact, nil
}

func (s service) Get(ctx context.Context) (Contact, error) {

	contact := Contact{}

	return contact, nil
}

func (s service) GetAll(ctx context.Context) ([]Contact, error) {

	contact := []Contact{}

	return contact, nil
}
