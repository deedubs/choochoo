package handlers

import (
	"fmt"
	"net/http"
)

// HandleRoot provides information about the server endpoints
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Choochoo GitHub Webhook Server\nEndpoints:\n- POST /webhook - GitHub webhook endpoint\n- GET /health - Health check\n")
}