package postal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// WebhookEvent represents an event from Postal
type WebhookEvent struct {
	Event     string                 `json:"event"`
	Timestamp time.Time              `json:"timestamp"`
	UUID      string                 `json:"uuid"`
	Payload   map[string]interface{} `json:"payload"`
}

// Common webhook event types
const (
	EventMessageSent      = "MessageSent"
	EventMessageDelivered = "MessageDelivered"
	EventMessageBounced   = "MessageBounced"
	EventMessageFailed    = "MessageFailed"
	EventMessageHeld      = "MessageHeld"
	EventMessageOpened    = "MessageOpened"
	EventMessageClicked   = "MessageLinkClicked"
)

// VerifyWebhookSignature verifies the HMAC signature of a webhook payload
func VerifyWebhookSignature(payload []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// ParseWebhookEvent parses a webhook event from JSON
func ParseWebhookEvent(data []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook event: %w", err)
	}
	return &event, nil
}

// GetMessageIDFromPayload extracts the message ID from webhook payload
func GetMessageIDFromPayload(payload map[string]interface{}) string {
	if msgID, ok := payload["message_id"].(string); ok {
		return msgID
	}
	if msg, ok := payload["message"].(map[string]interface{}); ok {
		if msgID, ok := msg["id"].(string); ok {
			return msgID
		}
	}
	return ""
}
