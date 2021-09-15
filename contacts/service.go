package contacts

import (
	"context"
	"errors"
	"fmt"
	"github.com/ncostamagna/axul_contact/pkg/client"
	"github.com/ncostamagna/streetflow/slack"
	"github.com/ncostamagna/streetflow/telegram"
	"strconv"
	"time"

	"github.com/ncostamagna/rerrors"

	"github.com/go-kit/kit/log"
)

//Service interface
type Service interface {
	Create(ctx context.Context, contact *Contact) rerrors.RestErr
	Update(ctx context.Context) (*Contact, rerrors.RestErr)
	Get(ctx context.Context, id string) (*Contact, rerrors.RestErr)
	GetAll(ctx context.Context, contacts *[]Contact, f Filter) rerrors.RestErr
	Alert(ctx context.Context, contacts *[]Contact, birthday string) rerrors.RestErr
	authorization(ctx context.Context, id, token string) error
}

type service struct {
	repo      Repository
	slackTran *slack.SlackBuilder
	telegTran *telegram.Transport
	tempTran  Transport
	userTran  client.Transport
	logger    log.Logger
}

type updateCb func(uint, time.Time) error

//NewService is a service handler
func NewService(repo Repository, slackTran *slack.SlackBuilder, telegTran *telegram.Transport, tempTran Transport, userTran client.Transport, logger log.Logger) Service {
	return &service{
		repo:      repo,
		slackTran: slackTran,
		telegTran: telegTran,
		userTran: userTran,
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

func (s service) Get(ctx context.Context, id string) (*Contact, rerrors.RestErr) {

	contact := Contact{}
	if err := s.repo.Get(ctx, &contact, id); err != nil {
		return nil, rerrors.NewBadRequestError(err)
	}

	return &contact, nil
}

func (s service) GetAll(ctx context.Context, contacts *[]Contact, f Filter) rerrors.RestErr {

	days, err := strconv.Atoi(f.birthday)

	if err == nil {
		if err := s.repo.GetByBirthdayRange(ctx, contacts, days); err != nil {
			return rerrors.NewInternalServerError(err)
		}
		return nil
	}

	if err := s.repo.GetAll(ctx, contacts, f); err != nil {
		return rerrors.NewInternalServerError(err)
	}

	return nil
}

func (s service) Alert(ctx context.Context, contacts *[]Contact, birthday string) rerrors.RestErr {

	days, err := strconv.Atoi(birthday)

	if err != nil {
		days = 0
	}

	if err := s.repo.GetByBirthdayRange(ctx, contacts, days); err != nil {
		return rerrors.NewInternalServerError(err)
	}

	for _, contact := range *contacts {
		fmt.Println(contact)
		switch days {
		case 1, 3:
			//slack alert
			fmt.Println("Slack Alert")
			res := s.slackTran.SendMessage("<@U01CDEPA3T9> " + message(days, contact.Nickname, contact.Phone))
			fmt.Println(res)
		case 0:
			//telegra alert
			fmt.Println("Telegram Alert")
			err := telegram.NewTelegramBuilder(*s.telegTran).Message(message(days, contact.Nickname, contact.Phone)).Send()
			fmt.Println(err)
		}
	}

	return nil
}

func message(days int, nickname, phone string) string {

	switch days {
	case 1:
		return fmt.Sprintf("Mañana es el cumpleaños de %s, recorda saludarlo", nickname)
	case 3:
		return fmt.Sprintf("En 3 dias es el cumpleaños de %s, recorda saludarlo", nickname)
	case 0:
		return "Hola Nahuel,\nhoy es el cumple de " + nickname + ", recorda saludarlo en su dia\n\nhttps://wa.me/" + phone + "?text=Feliz%20cumple%20" + nickname + "%0AEspero%20que%20lo%20pases%20de%20lo%20mejor!%0ATe%20mando%20un%20abrazo%20y%20muchos%20exitos!"
	}

	return ""
}

func (s *service) authorization(ctx context.Context, id, token string) error {
	a, err := s.userTran.GetAuth(id, token)

	if err != nil {
		fmt.Println(err)
		return errors.New("invalid authentication")
	}

	if a < 1 {
		return errors.New("invalid authorization")
	}

	return nil
}
