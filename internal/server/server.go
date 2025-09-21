package server

import (
	"log"
	"net/http"
	"os"

	"github.com/deedubs/choochoo/internal/handlers"
)

// WebhookServer represents the main server
type WebhookServer struct {
	webhookSecret string
	port          string
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

// Start starts the webhook server
func (ws *WebhookServer) Start() {
	mux := http.NewServeMux()
	
	// Create handlers with the webhook secret for signature validation
	webhookHandler := handlers.NewWebhookHandler(ws.webhookSecret)
	healthHandler := handlers.NewHealthHandler()
	
	// Register routes
	mux.HandleFunc("/webhook", webhookHandler.HandleWebhook)
	mux.HandleFunc("/health", healthHandler.HandleHealth)
	mux.HandleFunc("/", handlers.HandleRoot)

	log.Printf("Starting choochoo webhook server on port %s", ws.port)
	log.Printf("Webhook endpoint: http://localhost:%s/webhook", ws.port)
	log.Printf("Health check: http://localhost:%s/health", ws.port)
	
	if err := http.ListenAndServe(":"+ws.port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}