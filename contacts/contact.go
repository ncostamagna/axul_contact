package contacts

import (
	"context"
	"time"

	"gorm.io/gorm"
)

//Contact model
type Contact struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	Firstname  string         `gorm:"size:30" json:"firstname"`
	Lastname   string         `gorm:"size:30" json:"lastname"`
	Nickname   string         `gorm:"size:30" json:"nickname"`
	Gender     string         `gorm:"size:1" json:"gender"`
	Phone      string         `gorm:"size:30" json:"phone"`
	Photo      string         `gorm:"size:200" json:"photo"`
	Birthday   time.Time      `json:"birthday"`
	TemplateID uint           `json:"template_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

//Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	Update(ctx context.Context, contact *Contact, contactValues Contact) error
	GetAll(ctx context.Context, contact *[]Contact) error
	Get(ctx context.Context, contact *Contact, id uint) error
	GetByBirthdayRange(ctx context.Context, contacts *[]Contact, days int) error
}
