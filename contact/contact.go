package emails

import (
	"context"
	"time"

	"gorm.io/gorm"
)

//Contact model
type Contact struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `json:"-"`
	FirstName  string `gorm:"size:70" json:"subject"`
	LastName   string `gorm:"size:50" json:"contentType"`
	NickName   string `gorm:"type:text" json:"contentValue"`
	Gender     string `json:"sandboxMode"` //`gorm:"size:50"`
	Phone      string `json:"-"`
	Photo      string
	BirthDay   time.Time
	TemplateID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

//Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	Update(ctx context.Context, contact *Contact, contactValues Contact) error
	GetAll(ctx context.Context, contact *[]Contact) error
	Get(ctx context.Context, contact *Contact, id uint) error
}
