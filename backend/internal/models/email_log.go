package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailStatus string

const (
	EmailStatusQueued    EmailStatus = "queued"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusDelivered EmailStatus = "delivered"
	EmailStatusBounced   EmailStatus = "bounced"
	EmailStatusFailed    EmailStatus = "failed"
	EmailStatusComplaint EmailStatus = "complaint"
)

type EmailLog struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	MessageID       string         `gorm:"uniqueIndex;not null" json:"message_id"`
	PostalMessageID string         `gorm:"index" json:"postal_message_id"`
	From            string         `gorm:"not null" json:"from"`
	To              string         `gorm:"not null;index" json:"to"`
	Subject         string         `json:"subject"`
	Status          EmailStatus    `gorm:"default:'queued';index" json:"status"`
	ErrorMessage    string         `json:"error_message,omitempty"`
	OpenedAt        *time.Time     `json:"opened_at,omitempty"`
	ClickedAt       *time.Time     `json:"clicked_at,omitempty"`
	BouncedAt       *time.Time     `json:"bounced_at,omitempty"`
	DeliveredAt     *time.Time     `json:"delivered_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}
