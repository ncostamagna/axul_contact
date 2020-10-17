package contacts

import (
	"context"

	"github.com/go-kit/kit/log"
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

	return nil

}

func (repo *repo) GetAll(ctx context.Context, contact *[]Contact) error {

	return nil
}

func (repo *repo) Get(ctx context.Context, contact *Contact, id uint) error {

	return nil
}

func (repo *repo) Update(ctx context.Context, contact *Contact, contactValues Contact) error {

	return nil
}
