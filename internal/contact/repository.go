package contact

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/google/uuid"
	"github.com/ncostamagna/axul_domain/domain"
	"gorm.io/gorm"
)

// Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, contact *domain.Contact) error
	Update(ctx context.Context, contact *domain.Contact, contactValues domain.Contact) error
	GetAll(ctx context.Context, f Filter, offset, limit int) ([]domain.Contact, error)
	Get(ctx context.Context, id string) (*domain.Contact, error)
	Count(ctx context.Context, filters Filter) (int, error)
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

func (repo *repo) Create(ctx context.Context, contact *domain.Contact) error {

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

func (repo *repo) GetAll(ctx context.Context, f Filter, offset, limit int) ([]domain.Contact, error) {

	var tx *gorm.DB
	var cs []domain.Contact

	tx = repo.db.Model(&cs)
	tx = applyFilters(tx, f)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Find(&cs)

	for i := range cs {
		year := f.firstDate.Year()
		if cs[i].Birthday.Month() < f.firstDate.Month() {
			year++
		} else if cs[i].Birthday.Month() == f.firstDate.Month() {
			if cs[i].Birthday.Day() < f.firstDate.Day() {
				year++
			}
		}

		bd := time.Date(year, cs[i].Birthday.Month(), cs[i].Birthday.Day(), 0, 0, 0, 0, time.UTC)
		cs[i].Days = int64(bd.Sub(f.firstDate).Hours() / 24)
	}

	sort.SliceStable(cs, func(i, j int) bool {
		return cs[i].Days < cs[j].Days
	})

	if result.Error != nil {
		return nil, repo.log.CatchError(result.Error)
	}

	_ = repo.log.CatchMessage(fmt.Sprintf("Row: %d", result.RowsAffected))

	return cs, nil
}

func (repo *repo) Get(_ context.Context, id string) (*domain.Contact, error) {
	contact := domain.Contact{}

	if err := repo.db.Where("id = ?", id).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (repo *repo) Update(ctx context.Context, contact *domain.Contact, contactValues domain.Contact) error {

	return nil
}

func (repo *repo) Delete(ctx context.Context, contact *[]domain.Contact) error {

	return nil
}

func (r *repo) Count(ctx context.Context, filters Filter) (int, error) {
	var count int64

	tx := r.db.WithContext(ctx).Model(domain.Contact{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, f Filter) *gorm.DB {

	if f.RangeDays != nil {
		second := f.firstDate.AddDate(0, 0, int(*f.RangeDays)).Add(time.Hour * 20)
		tx = tx.Where("CONCAT('"+strconv.Itoa(f.firstDate.Year())+"',DATE_FORMAT(birthday,'%m%d')) between DATE_FORMAT(?,'%Y%m%d') and DATE_FORMAT(?,'%Y%m%d')", f.firstDate, second)
	}

	if f.Name != "" {
		tx = tx.Where("UPPER(CONCAT(firstname, ' ', lastname, ' ', nickname)) like CONCAT('%',UPPER(?),'%')", f.Name)
	}

	if f.Month != 0 {
		tx = tx.Where("MONTH(birthday) = ?", f.Month)
	}

	if f.Birthday != nil {
		date := time.Now().AddDate(0, 0, *f.Birthday)
		day, month := date.Day(), int(date.Month())
		tx = tx.Where("month(birthday) = ? and day(birthday) = ?", month, day)
	}

	return tx
}
