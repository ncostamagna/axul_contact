package contact

import (
	"context"
	"fmt"
	"github.com/digitalhouse-dev/dh-kit/logger"
	authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_domain/domain"
	"github.com/starry-axul/notifit-go-sdk/notify"
	"strconv"
	"time"
	"os"
)

// Service interface
type Service interface {
	Create(ctx context.Context, firstName, lastName, nickName, gender, phone string, birthday time.Time) (*domain.Contact, error)
	Update(ctx context.Context, id string, firstName, lastName, nickName, gender, phone *string, birthday *time.Time) error
	Get(ctx context.Context, id string) (*domain.Contact, error)
	GetAll(ctx context.Context, f Filter, offset, limit int) ([]domain.Contact, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filters Filter) (int, error)
	Alert(ctx context.Context, birthday string) ([]domain.Contact, error)
	//authorization(ctx context.Context, id, token string) error
}

type service struct {
	repo      Repository
	notif     notify.Transport
	auth      authentication.Auth
	logger    logger.Logger
}

type Filter struct {
	RangeDays *int64
	Birthday  *int
	Firstname string
	Lastname  string
	Month     int16
	firstDate time.Time
}

// NewService is a service handler
func NewService(repo Repository, notif notify.Transport, auth authentication.Auth, logger logger.Logger) Service {
	return &service{
		repo:      repo,
		auth:      auth,
		notif:     notif,
		logger:    logger,
	}
}

// Create service
func (s service) Create(ctx context.Context, firstName, lastName, nickName, gender, phone string, birthday time.Time) (*domain.Contact, error) {

	c := domain.Contact{
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

func (s service) Update(ctx context.Context, id string, firstName, lastName, nickName, gender, phone *string, birthday *time.Time) error {
	return s.repo.Update(ctx, id, firstName, lastName, nickName, gender, phone, birthday)
}

func (s service) Delete(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Contact, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s service) GetAll(ctx context.Context, f Filter, offset, limit int) ([]domain.Contact, error) {

	cs, err := s.repo.GetAll(ctx, f, offset, limit)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (s service) Alert(ctx context.Context, birthday string) ([]domain.Contact, error) {

	days, err := strconv.Atoi(birthday)
	if err != nil {
		days = 0
	}

	cs, err := s.repo.GetAll(ctx, Filter{Birthday: &days}, 0, 0)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {

		if days == 0 {
			if err := s.notif.Push(ctx, fmt.Sprintf(os.Getenv("BIRTHDAY_TITLE"), c.Firstname, c.Lastname), fmt.Sprintf(os.Getenv("BIRTHDAY_TEXT"), c.Firstname, c.Lastname), os.Getenv("BIRTHDAY_PAGE")); err != nil {
				return nil, err
			}
		}
	}

	return cs, nil
}

/* deprecated
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
}*/

func (s service) Count(ctx context.Context, filters Filter) (int, error) {
	return s.repo.Count(ctx, filters)
}
