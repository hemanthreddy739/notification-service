package main

import (
	"encoding/json"
	"net/http"
)

// NotificationRequest represents a notification request payload
type NotificationRequest struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// NotificationResponse represents the response from sending a notification
type NotificationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendNotification handles the notification endpoint
func SendNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		})
		return
	}

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if req.Message == "" || req.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "message and user_id are required",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NotificationResponse{
		Status:  "success",
		Message: "notification sent",
	})
}

// HealthCheck handles the health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func main() {
	http.HandleFunc("/notify", SendNotification)
	http.HandleFunc("/health", HealthCheck)
	http.ListenAndServe(":8080", nil)
}
