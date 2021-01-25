package contacts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ncostamagna/streetflow/slack"
	"github.com/ncostamagna/streetflow/telegram"

	"github.com/ncostamagna/rerrors"

	"github.com/go-kit/kit/log"
)

type service struct {
	repo      Repository
	slackTran slack.SlackBuilder
	telegTran telegram.Transport
	tempTran  Transport
	logger    log.Logger
}

type updateCb func(uint, time.Time) error

//NewService is a service handler
func NewService(repo Repository, slackTran slack.SlackBuilder, telegTran telegram.Transport, logger log.Logger) Service {
	return &service{
		repo:      repo,
		slackTran: slackTran,
		telegTran: telegTran,
		logger:    logger,
	}
}

//Create service
func (s service) Create(ctx context.Context, contact *Contact) rerrors.RestErr {

	err := s.repo.Create(ctx, contact)

	if err != nil {
		return rerrors.NewInternalServerError(err)
	}

	return nil
}

func (s service) Update(ctx context.Context) (*Contact, rerrors.RestErr) {

	contact := Contact{}

	return &contact, nil
}

func (s service) Delete(ctx context.Context) (*Contact, rerrors.RestErr) {

	contact := Contact{}

	return &contact, nil
}

func (s service) Get(ctx context.Context) (Contact, rerrors.RestErr) {

	contact := Contact{}

	return contact, nil
}

func (s service) GetAll(ctx context.Context, contacts *[]Contact, birthday string) rerrors.RestErr {

	days, err := strconv.Atoi(birthday)

	fmt.Println(days)
	fmt.Println(err)
	if err == nil {
		if err := s.repo.GetByBirthdayRange(ctx, contacts, days); err != nil {
			return rerrors.NewInternalServerError(err)
		}
		return nil
	}

	if err := s.repo.GetAll(ctx, contacts); err != nil {
		return rerrors.NewInternalServerError(err)
	}

	return nil
}

func (s service) Alert(ctx context.Context, contacts *[]Contact, birthday string) rerrors.RestErr {

	days, err := strconv.Atoi(birthday)

	if err == nil {
		if err := s.repo.GetByBirthdayRange(ctx, contacts, days); err != nil {
			return rerrors.NewInternalServerError(err)
		}

		fmt.Println(contacts)
		switch days {
		case 1, 3:
			//slack alert
		case 0:
			//telegra alert
		}
		return nil
	}

	return nil
}
