package models

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Key         string         `gorm:"uniqueIndex;not null" json:"key"` // Hashed
	KeyPrefix   string         `gorm:"not null" json:"key_prefix"`      // First 8 chars for display
	Permissions string         `gorm:"type:jsonb" json:"permissions"`   // JSON array of permissions
	RateLimit   int            `gorm:"default:1000" json:"rate_limit"`  // Requests per hour
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}
