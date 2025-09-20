package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler_HandleHealth(t *testing.T) {
	handler := NewHealthHandler()
	
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	
	handler.HandleHealth(rr, req)
	
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