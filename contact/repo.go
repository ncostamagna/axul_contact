package emails

import (
	"context"
	"errors"
	"fmt"

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

func (repo *repo) Create(ctx context.Context, email *Email) error {

	logger := log.With(repo.logger, "method", "Create")

	result := repo.db.Create(&email)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)
	_ = logger.Log("ID", email.ID)

	return nil

}

func (repo *repo) GetAll(ctx context.Context, emails *[]Email, wasSent string) error {

	logger := log.With(repo.logger, "method", "GetAll")
	var w = repo.db

	switch wasSent {
	case "1":
		w = w.Where("not sent_at is NULL")
	case "0":
		w = w.Where("sent_at IS NULL")
	}

	result := w.Find(&emails)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	return nil

}

func (repo *repo) Get(ctx context.Context, email *Email, id uint) error {

	logger := log.With(repo.logger, "method", "Get")

	if id == 0 {
		return errors.New("Invalid ID value")
	}
	result := repo.db.Where(&Email{ID: id}).Preload("Addresses").Preload("Attachment").First(&email)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	fmt.Println(email)
	return nil

}

func (repo *repo) Update(ctx context.Context, model *Email, emailValues Email) error {

	logger := log.With(repo.logger, "method", "Update")
	result := repo.db.Model(&model).UpdateColumns(emailValues)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)
	_ = logger.Log("ID", model.ID)

	return nil

}
