package contact

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ncostamagna/axul_contact/pkg/client"
	"github.com/ncostamagna/streetflow/slack"
	"github.com/ncostamagna/streetflow/telegram"

	"github.com/digitalhouse-dev/dh-kit/logger"
	authentication "github.com/ncostamagna/axul_auth/auth"
)

// Service interface
type Service interface {
	Create(ctx context.Context, firstName, lastName, nickName, gender, phone string, birthday time.Time) (*Contact, error)
	Update(ctx context.Context, id, firstName, lastName, nickName, gender, phone string, birthday time.Time) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Contact, error)
	GetAll(ctx context.Context, contacts *[]Contact, f Filter) error
	Alert(ctx context.Context, contacts *[]Contact, birthday string) error
	authorization(ctx context.Context, id, token string) error
}

type service struct {
	repo      Repository
	slackTran *slack.SlackBuilder
	telegTran *telegram.Transport
	userTran  client.Transport
	auth      authentication.Auth
	logger    logger.Logger
}

// NewService is a service handler
func NewService(repo Repository, slackTran *slack.SlackBuilder, telegTran *telegram.Transport, tempTran Transport, userTran client.Transport, auth authentication.Auth, logger logger.Logger) Service {
	return &service{
		repo:      repo,
		slackTran: slackTran,
		telegTran: telegTran,
		userTran:  userTran,
		auth:      auth,
		logger:    logger,
	}
}

// Create service
func (s service) Create(ctx context.Context, firstName, lastName, nickName, gender, phone string, birthday time.Time) (*Contact, error) {

	c := Contact{
		Firstname: firstName,
		Lastname:  lastName,
		Nickname:  nickName,
		Gender:    gender,
		Phone:     phone,
		Birthday:  birthday,
	}

	if err := s.repo.Create(ctx, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (s service) Update(ctx context.Context, id, firstName, lastName, nickName, gender, phone string, birthday time.Time) error {
	return nil
}

func (s service) Delete(ctx context.Context, id string) error {
	return nil
}

func (s service) Get(ctx context.Context, id string) (*Contact, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s service) GetAll(ctx context.Context, contacts *[]Contact, f Filter) error {

	days, err := strconv.Atoi(f.birthday)

	if err == nil {
		if err := s.repo.GetByBirthdayRange(ctx, contacts, days); err != nil {
			return err
		}
		return nil
	}

	if err := s.repo.GetAll(ctx, contacts, f); err != nil {
		return err
	}

	return nil
}

func (s service) Alert(ctx context.Context, contacts *[]Contact, birthday string) error {

	days, err := strconv.Atoi(birthday)

	if err != nil {
		days = 0
	}

	if err := s.repo.GetByBirthdayRange(ctx, contacts, days); err != nil {
		return err
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
	fmt.Println(id, token)
	return s.auth.Access(id, token)
}
