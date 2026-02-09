package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shohag/seentics-email/internal/database"
	"github.com/shohag/seentics-email/internal/models"
	"github.com/shohag/seentics-email/internal/postal"
)

type EmailHandler struct {
	postalClient *postal.Client
}

func NewEmailHandler(postalClient *postal.Client) *EmailHandler {
	return &EmailHandler{
		postalClient: postalClient,
	}
}

type SendEmailRequest struct {
	To        []string          `json:"to" binding:"required"`
	From      string            `json:"from" binding:"required,email"`
	Subject   string            `json:"subject" binding:"required"`
	HTMLBody  string            `json:"html_body"`
	PlainBody string            `json:"plain_body"`
	Headers   map[string]string `json:"headers"`
}

type SendEmailResponse struct {
	MessageID       string `json:"message_id"`
	PostalMessageID string `json:"postal_message_id"`
	Status          string `json:"status"`
}

// SendEmail sends an email via Postal
func (h *EmailHandler) SendEmail(c *gin.Context) {
	userID := c.GetUint("userID")

	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that at least one body is provided
	if req.HTMLBody == "" && req.PlainBody == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either html_body or plain_body is required"})
		return
	}

	// Generate unique message ID
	messageID := uuid.New().String()

	// Send via Postal
	postalReq := postal.SendEmailRequest{
		To:        req.To,
		From:      req.From,
		Subject:   req.Subject,
		HTMLBody:  req.HTMLBody,
		PlainBody: req.PlainBody,
		Headers:   req.Headers,
	}

	postalResp, err := h.postalClient.SendEmail(postalReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send email: %v", err)})
		return
	}

	// Extract Postal message ID
	var postalMessageID string
	if postalResp.Data.MessageID != "" {
		postalMessageID = postalResp.Data.MessageID
	}

	// Log email in database
	for _, recipient := range req.To {
		emailLog := models.EmailLog{
			UserID:          userID,
			MessageID:       messageID,
			PostalMessageID: postalMessageID,
			From:            req.From,
			To:              recipient,
			Subject:         req.Subject,
			Status:          models.EmailStatusSent,
		}

		if err := database.DB.Create(&emailLog).Error; err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to log email: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, SendEmailResponse{
		MessageID:       messageID,
		PostalMessageID: postalMessageID,
		Status:          "sent",
	})
}

// ListEmails returns paginated email logs
func (h *EmailHandler) ListEmails(c *gin.Context) {
	userID := c.GetUint("userID")

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Filter parameters
	status := c.Query("status")
	recipient := c.Query("to")

	query := database.DB.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if recipient != "" {
		query = query.Where("to LIKE ?", "%"+recipient+"%")
	}

	// Get total count
	var total int64
	query.Model(&models.EmailLog{}).Count(&total)

	// Get emails
	var emails []models.EmailLog
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&emails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"emails": emails,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetEmail returns details of a specific email
func (h *EmailHandler) GetEmail(c *gin.Context) {
	userID := c.GetUint("userID")
	emailID := c.Param("id")

	var email models.EmailLog
	if err := database.DB.Where("id = ? AND user_id = ?", emailID, userID).First(&email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	// Optionally fetch latest status from Postal
	if email.PostalMessageID != "" && h.postalClient != nil {
		if details, err := h.postalClient.GetMessage(email.PostalMessageID); err == nil {
			// Update status based on Postal response
			email.Status = models.EmailStatus(details.Status)
		}
	}

	c.JSON(http.StatusOK, email)
}
