package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/deedubs/choochoo/internal/webhook"
)

// WebhookHandler handles GitHub webhook requests
type WebhookHandler struct {
	webhookSecret string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(secret string) *WebhookHandler {
	return &WebhookHandler{
		webhookSecret: secret,
	}
}

// validateSignature validates the GitHub webhook signature
func (wh *WebhookHandler) validateSignature(payload []byte, signature string) bool {
	if wh.webhookSecret == "" {
		return true // Skip validation if no secret is set
	}

	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	// Remove "sha256=" prefix
	providedSignature := signature[7:]

	// Compute the expected signature
	mac := hmac.New(sha256.New, []byte(wh.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures using hmac.Equal for constant-time comparison
	providedBytes, err := hex.DecodeString(providedSignature)
	if err != nil {
		return false
	}
	expectedBytes, err := hex.DecodeString(expectedSignature)
	if err != nil {
		return false
	}

	return hmac.Equal(providedBytes, expectedBytes)
}

// HandleWebhook processes incoming GitHub webhook requests
func (wh *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get GitHub headers
	eventType := r.Header.Get("X-GitHub-Event")
	deliveryID := r.Header.Get("X-GitHub-Delivery")
	signature := r.Header.Get("X-Hub-Signature-256")

	// Validate signature if webhook secret is configured
	if !wh.validateSignature(body, signature) {
		log.Printf("Invalid signature for delivery %s", deliveryID)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse the JSON payload
	var event webhook.GitHubEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing JSON payload: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Log the webhook event
	repoName := "unknown"
	if event.Repository != nil {
		if name, ok := event.Repository["full_name"].(string); ok {
			repoName = name
		}
	}

	senderLogin := "unknown"
	if event.Sender != nil {
		if login, ok := event.Sender["login"].(string); ok {
			senderLogin = login
		}
	}

	log.Printf("Received %s event from %s (delivery: %s, sender: %s)", 
		eventType, repoName, deliveryID, senderLogin)

	if event.Action != "" {
		log.Printf("Event action: %s", event.Action)
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": "Webhook received and processed",
	}
	json.NewEncoder(w).Encode(response)
}