package contacts

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"

	"github.com/jinzhu/gorm"
)

//Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	Update(ctx context.Context, contact *Contact, contactValues Contact) error
	GetAll(ctx context.Context, contact *[]Contact) error
	Get(ctx context.Context, contact *Contact, id string) error
	GetByBirthdayRange(ctx context.Context, contacts *[]Contact, days int) error
}

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

//NewRepo is a repositories handler
func NewRepo(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) Create(ctx context.Context, contact *Contact) error {

	logger := log.With(repo.logger, "method", "Create")

	contact.ID = uuid.New().String()
	result := repo.db.Create(&contact)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)
	_ = logger.Log("ID", contact.ID)

	return nil
}

func (repo *repo) GetAll(ctx context.Context, contact *[]Contact) error {

	logger := log.With(repo.logger, "method", "GetAll")

	result := repo.db.Find(&contact)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)

	return nil
}

func (repo *repo) Get(ctx context.Context, contact *Contact, id string) error {
	result := repo.db.Where("id = ?", id).First(&contact)
	return result.Error
}

func (repo *repo) GetByBirthdayRange(ctx context.Context, contacts *[]Contact, days int) error {

	date := time.Now().AddDate(0, 0, days)
	day, month := date.Day(), int(date.Month())
	repo.db.Where("month(birthday) = ? and day(birthday) = ?", month, day).Find(&contacts)
	return nil
}

func (repo *repo) Update(ctx context.Context, contact *Contact, contactValues Contact) error {

	return nil
}

func (repo *repo) Delete(ctx context.Context, contact *[]Contact) error {

	return nil
}
