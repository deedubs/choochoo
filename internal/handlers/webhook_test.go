package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test helper functions

func generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// Tests for WebhookHandler

func TestWebhookHandler_ValidateSignature_NoSecret(t *testing.T) {
	handler := NewWebhookHandler("", nil)
	payload := []byte(`{"test": "data"}`)
	
	// Should return true when no secret is set (skip validation)
	result := handler.validateSignature(payload, "any-signature")
	if !result {
		t.Error("Expected validation to pass when no secret is set")
	}
}

func TestWebhookHandler_ValidateSignature_ValidSignature(t *testing.T) {
	secret := "test-secret"
	handler := NewWebhookHandler(secret, nil)
	payload := []byte(`{"test": "data"}`)
	signature := generateSignature(payload, secret)
	
	result := handler.validateSignature(payload, signature)
	if !result {
		t.Error("Expected validation to pass with valid signature")
	}
}

func TestWebhookHandler_ValidateSignature_InvalidSignature(t *testing.T) {
	handler := NewWebhookHandler("test-secret", nil)
	payload := []byte(`{"test": "data"}`)
	
	result := handler.validateSignature(payload, "sha256=invalid-signature")
	if result {
		t.Error("Expected validation to fail with invalid signature")
	}
}

func TestWebhookHandler_ValidateSignature_MissingPrefix(t *testing.T) {
	handler := NewWebhookHandler("test-secret", nil)
	payload := []byte(`{"test": "data"}`)
	
	result := handler.validateSignature(payload, "invalid-without-prefix")
	if result {
		t.Error("Expected validation to fail with missing sha256= prefix")
	}
}

func TestWebhookHandler_ValidateSignature_InvalidHex(t *testing.T) {
	handler := NewWebhookHandler("test-secret", nil)
	payload := []byte(`{"test": "data"}`)
	
	result := handler.validateSignature(payload, "sha256=invalid-hex-data")
	if result {
		t.Error("Expected validation to fail with invalid hex data")
	}
}

func TestWebhookHandler_HandleWebhook_InvalidMethod(t *testing.T) {
	handler := NewWebhookHandler("", nil)
	
	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestWebhookHandler_HandleWebhook_ValidRequest_NoSecret(t *testing.T) {
	handler := NewWebhookHandler("", nil)
	
	payload := `{"action":"push","repository":{"full_name":"test/repo"},"sender":{"login":"testuser"}}`
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}
	
	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got %s", response["status"])
	}
}

func TestWebhookHandler_HandleWebhook_ValidRequest_WithSecret(t *testing.T) {
	secret := "test-secret"
	handler := NewWebhookHandler(secret, nil)
	
	payload := `{"action":"push","repository":{"full_name":"test/repo"},"sender":{"login":"testuser"}}`
	payloadBytes := []byte(payload)
	signature := generateSignature(payloadBytes, secret)
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", signature)
	
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}
	
	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got %s", response["status"])
	}
}

func TestWebhookHandler_HandleWebhook_InvalidSignature(t *testing.T) {
	handler := NewWebhookHandler("test-secret", nil)
	
	payload := `{"action":"push","repository":{"full_name":"test/repo"},"sender":{"login":"testuser"}}`
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid-signature")
	
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, status)
	}
}

func TestWebhookHandler_HandleWebhook_InvalidJSON(t *testing.T) {
	handler := NewWebhookHandler("", nil)
	
	payload := `invalid json`
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestWebhookHandler_HandleWebhook_EmptyPayload(t *testing.T) {
	handler := NewWebhookHandler("", nil)
	
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestWebhookHandler_HandleWebhook_GitHubEvent_OptionalFields(t *testing.T) {
	handler := NewWebhookHandler("", nil)
	
	payload := `{}`  // Empty payload with no optional fields
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "ping")
	req.Header.Set("X-GitHub-Delivery", "test-delivery-id")
	
	rr := httptest.NewRecorder()
	
	handler.HandleWebhook(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
}