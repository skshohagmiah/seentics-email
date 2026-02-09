package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	APIKeys   []APIKey   `gorm:"foreignKey:UserID" json:"-"`
	Domains   []Domain   `gorm:"foreignKey:UserID" json:"-"`
	EmailLogs []EmailLog `gorm:"foreignKey:UserID" json:"-"`
	Webhooks  []Webhook  `gorm:"foreignKey:UserID" json:"-"`
}
