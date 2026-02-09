package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shohag/seentics-email/internal/database"
	"github.com/shohag/seentics-email/internal/models"
)

type DomainHandler struct{}

func NewDomainHandler() *DomainHandler {
	return &DomainHandler{}
}

type AddDomainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

type DomainResponse struct {
	ID                 uint                `json:"id"`
	Domain             string              `json:"domain"`
	VerificationStatus models.DomainStatus `json:"verification_status"`
	DNSRecords         []DNSRecord         `json:"dns_records"`
	CreatedAt          string              `json:"created_at"`
}

type DNSRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority,omitempty"`
}

// ListDomains returns all domains for the authenticated user
func (h *DomainHandler) ListDomains(c *gin.Context) {
	userID := c.GetUint("userID")

	var domains []models.Domain
	if err := database.DB.Where("user_id = ?", userID).Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains"})
		return
	}

	response := make([]DomainResponse, len(domains))
	for i, domain := range domains {
		response[i] = DomainResponse{
			ID:                 domain.ID,
			Domain:             domain.Domain,
			VerificationStatus: domain.VerificationStatus,
			DNSRecords:         generateDNSRecords(domain.Domain),
			CreatedAt:          domain.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// AddDomain adds a new domain
func (h *DomainHandler) AddDomain(c *gin.Context) {
	userID := c.GetUint("userID")

	var req AddDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if domain already exists
	var existing models.Domain
	if err := database.DB.Where("domain = ?", req.Domain).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Domain already exists"})
		return
	}

	domain := models.Domain{
		UserID:             userID,
		Domain:             req.Domain,
		VerificationStatus: models.DomainStatusPending,
	}

	if err := database.DB.Create(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add domain"})
		return
	}

	c.JSON(http.StatusCreated, DomainResponse{
		ID:                 domain.ID,
		Domain:             domain.Domain,
		VerificationStatus: domain.VerificationStatus,
		DNSRecords:         generateDNSRecords(domain.Domain),
		CreatedAt:          domain.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetDomainVerification returns DNS records needed for verification
func (h *DomainHandler) GetDomainVerification(c *gin.Context) {
	userID := c.GetUint("userID")
	domainID := c.Param("id")

	var domain models.Domain
	if err := database.DB.Where("id = ? AND user_id = ?", domainID, userID).First(&domain).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":      domain.Domain,
		"dns_records": generateDNSRecords(domain.Domain),
		"status":      domain.VerificationStatus,
	})
}

// DeleteDomain removes a domain
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	userID := c.GetUint("userID")
	domainID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", domainID, userID).Delete(&models.Domain{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Domain deleted successfully"})
}

// VerifyDomain checks DNS records and updates verification status
func (h *DomainHandler) VerifyDomain(c *gin.Context) {
	userID := c.GetUint("userID")
	domainID := c.Param("id")

	var domain models.Domain
	if err := database.DB.Where("id = ? AND user_id = ?", domainID, userID).First(&domain).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// TODO: Implement actual DNS verification logic
	// For now, just return the current status
	c.JSON(http.StatusOK, gin.H{
		"domain":  domain.Domain,
		"status":  domain.VerificationStatus,
		"message": "DNS verification is not yet implemented. Please verify manually in Postal.",
	})
}

// generateDNSRecords generates the required DNS records for a domain
func generateDNSRecords(domain string) []DNSRecord {
	return []DNSRecord{
		{
			Type:     "MX",
			Name:     domain,
			Value:    "postal.yourdomain.com",
			Priority: 10,
		},
		{
			Type:  "TXT",
			Name:  domain,
			Value: "v=spf1 include:postal.yourdomain.com ~all",
		},
		{
			Type:  "TXT",
			Name:  "_dmarc." + domain,
			Value: "v=DMARC1; p=none; rua=mailto:dmarc@" + domain,
		},
		{
			Type:  "CNAME",
			Name:  "postal._domainkey." + domain,
			Value: "postal._domainkey.postal.yourdomain.com",
		},
	}
}
