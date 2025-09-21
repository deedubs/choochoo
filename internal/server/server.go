package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/deedubs/choochoo/internal/database"
	"github.com/deedubs/choochoo/internal/handlers"
)

// WebhookServer represents the main server
type WebhookServer struct {
	webhookSecret string
	port          string
	dbConn        *database.Connection
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

	// Initialize database connection if DATABASE_URL is set
	var dbConn *database.Connection
	if os.Getenv("DATABASE_URL") != "" {
		ctx := context.Background()
		var err error
		dbConn, err = database.NewConnection(ctx)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v. Webhooks will be logged but not stored.", err)
		} else {
			log.Println("Successfully connected to database")
		}
	} else {
		log.Println("Warning: DATABASE_URL not set. Webhooks will be logged but not stored in database.")
	}

	return &WebhookServer{
		webhookSecret: webhookSecret,
		port:          port,
		dbConn:        dbConn,
	}
}

// Start starts the webhook server
func (ws *WebhookServer) Start() {
	mux := http.NewServeMux()
	
	// Create handlers with the webhook secret for signature validation and database connection
	webhookHandler := handlers.NewWebhookHandler(ws.webhookSecret, ws.dbConn)
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