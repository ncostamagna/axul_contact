package contacts

import (
	"time"

	"gorm.io/gorm"
)

//Contact model
type Contact struct {
	ID         string         `gorm:"size:40;primaryKey" json:"id"`
	UserID     string         `json:"user_id"`
	Firstname  string         `gorm:"size:30" json:"firstname"`
	Lastname   string         `gorm:"size:30" json:"lastname"`
	Nickname   string         `gorm:"size:30" json:"nickname"`
	Gender     string         `gorm:"size:1" json:"gender"`
	TypeID     string         `gorm:"size:40" json:"type"`
	Phone      string         `gorm:"size:30" json:"phone"`
	Instagram  string         `gorm:"size:40" json:"instagram"`
	Photo      string         `gorm:"size:200" json:"photo"`
	Birthday   time.Time      `json:"birthday"`
	TemplateID string         `json:"template_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
