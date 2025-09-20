package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Test helper functions

func generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func createTestWebhookServer(secret string) *WebhookServer {
	return &WebhookServer{
		webhookSecret: secret,
		port:          "8080",
	}
}

// Tests for NewWebhookServer

func TestNewWebhookServer_DefaultPort(t *testing.T) {
	// Clear any existing PORT env var
	oldPort := os.Getenv("PORT")
	os.Unsetenv("PORT")
	defer func() {
		if oldPort != "" {
			os.Setenv("PORT", oldPort)
		}
	}()

	server := NewWebhookServer()
	if server.port != "8080" {
		t.Errorf("Expected default port 8080, got %s", server.port)
	}
}

func TestNewWebhookServer_CustomPort(t *testing.T) {
	oldPort := os.Getenv("PORT")
	os.Setenv("PORT", "3000")
	defer func() {
		if oldPort != "" {
			os.Setenv("PORT", oldPort)
		} else {
			os.Unsetenv("PORT")
		}
	}()

	server := NewWebhookServer()
	if server.port != "3000" {
		t.Errorf("Expected port 3000, got %s", server.port)
	}
}

func TestNewWebhookServer_WithSecret(t *testing.T) {
	oldSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	os.Setenv("GITHUB_WEBHOOK_SECRET", "test-secret")
	defer func() {
		if oldSecret != "" {
			os.Setenv("GITHUB_WEBHOOK_SECRET", oldSecret)
		} else {
			os.Unsetenv("GITHUB_WEBHOOK_SECRET")
		}
	}()

	server := NewWebhookServer()
	if server.webhookSecret != "test-secret" {
		t.Errorf("Expected webhook secret 'test-secret', got %s", server.webhookSecret)
	}
}

// Tests for validateSignature

func TestValidateSignature_NoSecret(t *testing.T) {
	server := createTestWebhookServer("")
	payload := []byte(`{"test": "data"}`)
	
	// Should return true when no secret is set (skip validation)
	result := server.validateSignature(payload, "any-signature")
	if !result {
		t.Error("Expected validation to pass when no secret is set")
	}
}

func TestValidateSignature_ValidSignature(t *testing.T) {
	secret := "test-secret"
	server := createTestWebhookServer(secret)
	payload := []byte(`{"test": "data"}`)
	signature := generateSignature(payload, secret)
	
	result := server.validateSignature(payload, signature)
	if !result {
		t.Error("Expected validation to pass with valid signature")
	}
}

func TestValidateSignature_InvalidSignature(t *testing.T) {
	server := createTestWebhookServer("test-secret")
	payload := []byte(`{"test": "data"}`)
	
	result := server.validateSignature(payload, "sha256=invalid-signature")
	if result {
		t.Error("Expected validation to fail with invalid signature")
	}
}

func TestValidateSignature_MissingPrefix(t *testing.T) {
	server := createTestWebhookServer("test-secret")
	payload := []byte(`{"test": "data"}`)
	
	result := server.validateSignature(payload, "invalid-signature-without-prefix")
	if result {
		t.Error("Expected validation to fail with signature missing sha256= prefix")
	}
}

func TestValidateSignature_InvalidHex(t *testing.T) {
	server := createTestWebhookServer("test-secret")
	payload := []byte(`{"test": "data"}`)
	
	result := server.validateSignature(payload, "sha256=invalid-hex-characters!")
	if result {
		t.Error("Expected validation to fail with invalid hex characters")
	}
}

// Tests for handleHealth

func TestHandleHealth(t *testing.T) {
	server := createTestWebhookServer("")
	
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	
	server.handleHealth(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	expected := `{"service":"choochoo-webhook-server","status":"healthy"}`
	body := strings.TrimSpace(rr.Body.String())
	if body != expected {
		t.Errorf("Expected body %s, got %s", expected, body)
	}
	
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

// Tests for handleWebhook

func TestHandleWebhook_InvalidMethod(t *testing.T) {
	server := createTestWebhookServer("")
	
	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestHandleWebhook_ValidRequest_NoSecret(t *testing.T) {
	server := createTestWebhookServer("")
	
	payload := map[string]interface{}{
		"action": "push",
		"repository": map[string]interface{}{
			"full_name": "test/repo",
		},
		"sender": map[string]interface{}{
			"login": "testuser",
		},
	}
	
	jsonPayload, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}
	
	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got %s", response["status"])
	}
}

func TestHandleWebhook_ValidRequest_WithSecret(t *testing.T) {
	secret := "test-secret"
	server := createTestWebhookServer(secret)
	
	payload := map[string]interface{}{
		"action": "push",
		"repository": map[string]interface{}{
			"full_name": "test/repo",
		},
		"sender": map[string]interface{}{
			"login": "testuser",
		},
	}
	
	jsonPayload, _ := json.Marshal(payload)
	signature := generateSignature(jsonPayload, secret)
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", signature)
	
	rr := httptest.NewRecorder()
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestHandleWebhook_InvalidSignature(t *testing.T) {
	server := createTestWebhookServer("test-secret")
	
	payload := map[string]interface{}{
		"action": "push",
		"repository": map[string]interface{}{
			"full_name": "test/repo",
		},
	}
	
	jsonPayload, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid-signature")
	
	rr := httptest.NewRecorder()
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, status)
	}
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	server := createTestWebhookServer("")
	
	invalidJSON := []byte(`{"invalid": json}`)
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestHandleWebhook_EmptyPayload(t *testing.T) {
	server := createTestWebhookServer("")
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

// Tests for GitHubEvent parsing

func TestGitHubEvent_OptionalFields(t *testing.T) {
	server := createTestWebhookServer("")
	
	// Test with minimal payload (no action, repository, or sender)
	payload := map[string]interface{}{}
	
	jsonPayload, _ := json.Marshal(payload)
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "ping")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	server.handleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
}

// Integration tests

func TestWebhookServer_RoutingIntegration(t *testing.T) {
	server := createTestWebhookServer("")
	
	// Create a test server with the same routing as the real server
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", server.handleWebhook)
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Choochoo GitHub Webhook Server"))
	})
	
	testServer := httptest.NewServer(mux)
	defer testServer.Close()
	
	// Test health endpoint
	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for health endpoint, got %d", resp.StatusCode)
	}
	
	// Test root endpoint
	resp, err = http.Get(testServer.URL + "/")
	if err != nil {
		t.Fatalf("Failed to call root endpoint: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for root endpoint, got %d", resp.StatusCode)
	}
	
	// Test 404 for unknown path
	resp, err = http.Get(testServer.URL + "/unknown")
	if err != nil {
		t.Fatalf("Failed to call unknown endpoint: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for unknown endpoint, got %d", resp.StatusCode)
	}
}