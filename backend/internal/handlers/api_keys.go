package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shohag/seentics-email/internal/database"
	"github.com/shohag/seentics-email/internal/models"
)

type APIKeyHandler struct{}

func NewAPIKeyHandler() *APIKeyHandler {
	return &APIKeyHandler{}
}

type CreateAPIKeyRequest struct {
	Name      string `json:"name" binding:"required"`
	RateLimit int    `json:"rate_limit" binding:"required,min=1"`
}

type APIKeyResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Key        string  `json:"key,omitempty"` // Only returned on creation
	KeyPrefix  string  `json:"key_prefix"`
	RateLimit  int     `json:"rate_limit"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
}

// ListAPIKeys returns all API keys for the authenticated user
func (h *APIKeyHandler) ListAPIKeys(c *gin.Context) {
	userID := c.GetUint("userID")

	var keys []models.APIKey
	if err := database.DB.Where("user_id = ?", userID).Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API keys"})
		return
	}

	response := make([]APIKeyResponse, len(keys))
	for i, key := range keys {
		var lastUsed *string
		if key.LastUsedAt != nil {
			formatted := key.LastUsedAt.Format("2006-01-02T15:04:05Z")
			lastUsed = &formatted
		}

		response[i] = APIKeyResponse{
			ID:         key.ID,
			Name:       key.Name,
			KeyPrefix:  key.KeyPrefix,
			RateLimit:  key.RateLimit,
			LastUsedAt: lastUsed,
			CreatedAt:  key.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateAPIKey generates a new API key
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate random API key
	rawKey, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Hash the key for storage
	hashedKey := hashAPIKey(rawKey)
	keyPrefix := rawKey[:8] // Store first 8 chars for display

	apiKey := models.APIKey{
		UserID:      userID,
		Name:        req.Name,
		Key:         hashedKey,
		KeyPrefix:   keyPrefix,
		RateLimit:   req.RateLimit,
		Permissions: "[]", // Default empty permissions
	}

	if err := database.DB.Create(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, APIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       rawKey, // Return the raw key only on creation
		KeyPrefix: apiKey.KeyPrefix,
		RateLimit: apiKey.RateLimit,
		CreatedAt: apiKey.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// DeleteAPIKey revokes an API key
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	userID := c.GetUint("userID")
	keyID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", keyID, userID).Delete(&models.APIKey{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

// UpdateAPIKey updates an API key's settings
func (h *APIKeyHandler) UpdateAPIKey(c *gin.Context) {
	userID := c.GetUint("userID")
	keyID := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		RateLimit int    `json:"rate_limit" binding:"min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var apiKey models.APIKey
	if err := database.DB.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.RateLimit > 0 {
		updates["rate_limit"] = req.RateLimit
	}

	if err := database.DB.Model(&apiKey).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key updated successfully"})
}

// generateAPIKey creates a random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(bytes), nil
}

func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
