package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shohag/seentics-email/internal/database"
	"github.com/shohag/seentics-email/internal/models"
	"github.com/shohag/seentics-email/internal/postal"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

type CreateWebhookRequest struct {
	URL    string   `json:"url" binding:"required,url"`
	Events []string `json:"events" binding:"required"`
}

type WebhookResponse struct {
	ID              uint     `json:"id"`
	URL             string   `json:"url"`
	Events          []string `json:"events"`
	Secret          string   `json:"secret,omitempty"` // Only on creation
	IsActive        bool     `json:"is_active"`
	LastTriggeredAt *string  `json:"last_triggered_at"`
	CreatedAt       string   `json:"created_at"`
}

// ListWebhooks returns all webhooks for the authenticated user
func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	userID := c.GetUint("userID")

	var webhooks []models.Webhook
	if err := database.DB.Where("user_id = ?", userID).Find(&webhooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch webhooks"})
		return
	}

	response := make([]WebhookResponse, len(webhooks))
	for i, wh := range webhooks {
		var lastTriggered *string
		if wh.LastTriggeredAt != nil {
			formatted := wh.LastTriggeredAt.Format("2006-01-02T15:04:05Z")
			lastTriggered = &formatted
		}

		response[i] = WebhookResponse{
			ID:              wh.ID,
			URL:             wh.URL,
			Events:          []string{}, // Parse from JSON
			IsActive:        wh.IsActive,
			LastTriggeredAt: lastTriggered,
			CreatedAt:       wh.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateWebhook creates a new webhook endpoint
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate webhook secret
	secret, err := generateWebhookSecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate webhook secret"})
		return
	}

	webhook := models.Webhook{
		UserID:   userID,
		URL:      req.URL,
		Events:   "[]", // Store as JSON
		Secret:   secret,
		IsActive: true,
	}

	if err := database.DB.Create(&webhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create webhook"})
		return
	}

	c.JSON(http.StatusCreated, WebhookResponse{
		ID:        webhook.ID,
		URL:       webhook.URL,
		Events:    req.Events,
		Secret:    secret, // Return secret only on creation
		IsActive:  webhook.IsActive,
		CreatedAt: webhook.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// DeleteWebhook removes a webhook
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	userID := c.GetUint("userID")
	webhookID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", webhookID, userID).Delete(&models.Webhook{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete webhook"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook deleted successfully"})
}

// HandlePostalWebhook receives webhooks from Postal
func (h *WebhookHandler) HandlePostalWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Parse webhook event
	event, err := postal.ParseWebhookEvent(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	// Extract message ID from payload
	messageID := postal.GetMessageIDFromPayload(event.Payload)
	if messageID == "" {
		c.JSON(http.StatusOK, gin.H{"message": "No message ID in payload"})
		return
	}

	// Update email log based on event type
	var emailLog models.EmailLog
	if err := database.DB.Where("postal_message_id = ?", messageID).First(&emailLog).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Email log not found"})
		return
	}

	now := time.Now()
	updates := make(map[string]interface{})

	switch event.Event {
	case postal.EventMessageDelivered:
		updates["status"] = models.EmailStatusDelivered
		updates["delivered_at"] = now
	case postal.EventMessageBounced:
		updates["status"] = models.EmailStatusBounced
		updates["bounced_at"] = now
	case postal.EventMessageFailed:
		updates["status"] = models.EmailStatusFailed
	case postal.EventMessageOpened:
		updates["opened_at"] = now
	case postal.EventMessageClicked:
		updates["clicked_at"] = now
	}

	if len(updates) > 0 {
		database.DB.Model(&emailLog).Updates(updates)
	}

	// Forward to user webhooks
	var webhooks []models.Webhook
	database.DB.Where("user_id = ? AND is_active = ?", emailLog.UserID, true).Find(&webhooks)

	// TODO: Implement webhook forwarding to user endpoints

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}

func generateWebhookSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "whsec_" + hex.EncodeToString(bytes), nil
}
