package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// WebhookServer represents the main server
type WebhookServer struct {
	webhookSecret string
	port          string
}

// GitHubEvent represents a generic GitHub webhook event
type GitHubEvent struct {
	Action     string                 `json:"action,omitempty"`
	Repository map[string]interface{} `json:"repository,omitempty"`
	Sender     map[string]interface{} `json:"sender,omitempty"`
}

// NewWebhookServer creates a new webhook server instance
func NewWebhookServer() *WebhookServer {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("Warning: GITHUB_WEBHOOK_SECRET not set. Webhook signature validation will be skipped.")
	}

	return &WebhookServer{
		webhookSecret: webhookSecret,
		port:          port,
	}
}

// validateSignature validates the GitHub webhook signature
func (ws *WebhookServer) validateSignature(payload []byte, signature string) bool {
	if ws.webhookSecret == "" {
		return true // Skip validation if no secret is set
	}

	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	// Remove "sha256=" prefix
	providedSignature := signature[7:]

	// Compute the expected signature
	mac := hmac.New(sha256.New, []byte(ws.webhookSecret))
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

// handleWebhook processes incoming GitHub webhook requests
func (ws *WebhookServer) handleWebhook(w http.ResponseWriter, r *http.Request) {
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
	if !ws.validateSignature(body, signature) {
		log.Printf("Invalid signature for delivery %s", deliveryID)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse the JSON payload
	var event GitHubEvent
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

// handleHealth provides a health check endpoint
func (ws *WebhookServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "healthy",
		"service": "choochoo-webhook-server",
	}
	json.NewEncoder(w).Encode(response)
}

// Start starts the webhook server
func (ws *WebhookServer) Start() {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/webhook", ws.handleWebhook)
	mux.HandleFunc("/health", ws.handleHealth)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Choochoo GitHub Webhook Server\nEndpoints:\n- POST /webhook - GitHub webhook endpoint\n- GET /health - Health check\n")
	})

	log.Printf("Starting choochoo webhook server on port %s", ws.port)
	log.Printf("Webhook endpoint: http://localhost:%s/webhook", ws.port)
	log.Printf("Health check: http://localhost:%s/health", ws.port)

	if err := http.ListenAndServe(":"+ws.port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func main() {
	server := NewWebhookServer()
	server.Start()
}
