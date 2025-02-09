package health

import (
	"encoding/json"
	"net/http"
)

// HealthResponse defines the response structure
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles the /health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
