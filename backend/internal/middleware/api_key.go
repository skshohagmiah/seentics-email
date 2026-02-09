package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shohag/seentics-email/internal/database"
	"github.com/shohag/seentics-email/internal/models"
)

type APIKeyMiddleware struct {
	redisClient *redis.Client
}

func NewAPIKeyMiddleware(redisClient *redis.Client) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		redisClient: redisClient,
	}
}

func (m *APIKeyMiddleware) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Hash the API key
		hashedKey := hashAPIKey(apiKey)

		// Check if key exists in database
		var key models.APIKey
		if err := database.DB.Where("key = ?", hashedKey).First(&key).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Check rate limit
		if !m.checkRateLimit(c.Request.Context(), key.ID, key.RateLimit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		// Update last used timestamp
		now := time.Now()
		database.DB.Model(&key).Update("last_used_at", now)

		// Set user ID and API key ID in context
		c.Set("userID", key.UserID)
		c.Set("apiKeyID", key.ID)
		c.Next()
	}
}

func (m *APIKeyMiddleware) checkRateLimit(ctx context.Context, keyID uint, limit int) bool {
	if m.redisClient == nil {
		return true // Skip rate limiting if Redis is not configured
	}

	key := fmt.Sprintf("ratelimit:apikey:%d", keyID)

	// Get current count
	count, err := m.redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return true // Allow on error
	}

	if count >= limit {
		return false
	}

	// Increment counter
	pipe := m.redisClient.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour)
	_, err = pipe.Exec(ctx)

	return err == nil
}

func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
