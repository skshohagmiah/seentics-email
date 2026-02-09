package models

import (
	"time"

	"gorm.io/gorm"
)

type Webhook struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	URL             string         `gorm:"not null" json:"url"`
	Events          string         `gorm:"type:jsonb" json:"events"` // JSON array of event types
	Secret          string         `gorm:"not null" json:"secret"`   // For webhook signature verification
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	LastTriggeredAt *time.Time     `json:"last_triggered_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}
