package models

import (
	"time"

	"gorm.io/gorm"
)

type DomainStatus string

const (
	DomainStatusPending  DomainStatus = "pending"
	DomainStatusVerified DomainStatus = "verified"
	DomainStatusFailed   DomainStatus = "failed"
)

type Domain struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"not null;index" json:"user_id"`
	Domain             string         `gorm:"uniqueIndex;not null" json:"domain"`
	VerificationStatus DomainStatus   `gorm:"default:'pending'" json:"verification_status"`
	PostalServerID     string         `json:"postal_server_id"`              // Postal's server ID
	PostalOrganization string         `json:"postal_organization"`           // Postal's organization
	DNSRecords         string         `gorm:"type:jsonb" json:"dns_records"` // JSON array of DNS records
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}
