package contact

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Filter struct {
	days     int64
	birthday string
	name     string
	month    int16
}

// Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	Update(ctx context.Context, contact *Contact, contactValues Contact) error
	GetAll(ctx context.Context, contact *[]Contact, f Filter) error
	Get(ctx context.Context, id string) (*Contact, error)
	GetByBirthdayRange(ctx context.Context, contacts *[]Contact, days int) error
}

type repo struct {
	db  *gorm.DB
	log logger.Logger
}

// NewRepo is a repositories handler
func NewRepo(db *gorm.DB, logger logger.Logger) Repository {
	return &repo{
		db:  db,
		log: logger,
	}
}

func (repo *repo) Create(ctx context.Context, contact *Contact) error {

	contact.ID = uuid.New().String()
	result := repo.db.Create(&contact)

	if result.Error != nil {
		_ = repo.log.CatchError
		return result.Error
	}

	_ = repo.log.CatchMessage(fmt.Sprintf("Row: %d", result.RowsAffected))
	_ = repo.log.CatchMessage(contact.ID)

	return nil
}

func (repo *repo) GetAll(ctx context.Context, contact *[]Contact, f Filter) error {

	var tx *gorm.DB

	tx = repo.db.Model(&contact)
	currentTime := time.Now().UTC()
	first := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	if f.days != 0 {
		second := first.AddDate(0, 0, int(f.days)).Add(time.Hour * 20)
		tx = tx.Where("CONCAT('"+strconv.Itoa(first.Year())+"',DATE_FORMAT(birthday,'%m%d')) between DATE_FORMAT(?,'%Y%m%d') and DATE_FORMAT(?,'%Y%m%d')", first, second)
	}

	if f.name != "" {
		tx = tx.Where("UPPER(CONCAT(firstname, ' ', lastname, ' ', nickname)) like CONCAT('%',UPPER(?),'%')", f.name)
	}

	if f.month != 0 {
		tx = tx.Where("MONTH(birthday) = ?", f.month)
	}

	result := tx.Find(&contact)

	for i := range *contact {
		year := currentTime.Year()
		if (*contact)[i].Birthday.Month() < currentTime.Month() {
			year++
		} else if (*contact)[i].Birthday.Month() == currentTime.Month() {
			if (*contact)[i].Birthday.Day() < currentTime.Day() {
				year++
			}
		}

		bd := time.Date(year, (*contact)[i].Birthday.Month(), (*contact)[i].Birthday.Day(), 0, 0, 0, 0, time.UTC)
		(*contact)[i].Days = int64(bd.Sub(first).Hours() / 24)
	}

	sort.SliceStable(*contact, func(i, j int) bool {
		return (*contact)[i].Days < (*contact)[j].Days
	})

	if result.Error != nil {
		return repo.log.CatchError(result.Error)
	}

	_ = repo.log.CatchMessage(fmt.Sprintf("Row: %d", result.RowsAffected))

	return nil
}

func (repo *repo) Get(_ context.Context, id string) (*Contact, error) {
	contact := Contact{}

	if err := repo.db.Where("id = ?", id).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
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
