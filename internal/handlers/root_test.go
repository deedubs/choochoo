package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot_ValidPath(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	
	HandleRoot(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	expected := "Choochoo GitHub Webhook Server\nEndpoints:\n- POST /webhook - GitHub webhook endpoint\n- GET /health - Health check\n"
	body := rr.Body.String()
	if body != expected {
		t.Errorf("Expected body %s, got %s", expected, body)
	}
}

func TestHandleRoot_InvalidPath(t *testing.T) {
	req := httptest.NewRequest("GET", "/invalid", nil)
	rr := httptest.NewRecorder()
	
	HandleRoot(rr, req)
	
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}
}