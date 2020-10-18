package contacts

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/jinzhu/gorm"
)

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

func (repo *repo) Get(ctx context.Context, contact *Contact, id uint) error {

	return nil
}

func (repo *repo) GetByBirthdayRange(ctx context.Context, contacts *[]Contact, days int) error {

	date := time.Now().AddDate(0, 0, -1*days)
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
