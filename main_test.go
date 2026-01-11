package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSendNotificationSuccess tests successful notification sending
func TestSendNotificationSuccess(t *testing.T) {
	payload := NotificationRequest{
		Message: "Test notification",
		UserID:  "user123",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	var response NotificationResponse
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	if response.Message != "notification sent" {
		t.Errorf("Expected message 'notification sent', got '%s'", response.Message)
	}
}

// TestSendNotificationMethodNotAllowed tests non-POST requests
func TestSendNotificationMethodNotAllowed(t *testing.T) {
	tests := []string{"GET", "PUT", "DELETE", "PATCH"}

	for _, method := range tests {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/notify", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SendNotification)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 405, got %v", status)
			}

			var response map[string]string
			json.NewDecoder(rr.Body).Decode(&response)

			if response["error"] != "method not allowed" {
				t.Errorf("Expected error 'method not allowed', got '%s'", response["error"])
			}
		})
	}
}

// TestSendNotificationInvalidBody tests invalid JSON body
func TestSendNotificationInvalidBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", status)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["error"] != "invalid request body" {
		t.Errorf("Expected error 'invalid request body', got '%s'", response["error"])
	}
}

// TestSendNotificationMissingMessage tests missing message field
func TestSendNotificationMissingMessage(t *testing.T) {
	payload := NotificationRequest{
		UserID: "user123",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", status)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["error"] != "message and user_id are required" {
		t.Errorf("Expected validation error, got '%s'", response["error"])
	}
}

// TestSendNotificationMissingUserID tests missing user_id field
func TestSendNotificationMissingUserID(t *testing.T) {
	payload := NotificationRequest{
		Message: "Test notification",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", status)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["error"] != "message and user_id are required" {
		t.Errorf("Expected validation error, got '%s'", response["error"])
	}
}

// TestSendNotificationEmptyMessage tests empty message field
func TestSendNotificationEmptyMessage(t *testing.T) {
	payload := NotificationRequest{
		Message: "",
		UserID:  "user123",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", status)
	}
}

// TestSendNotificationEmptyUserID tests empty user_id field
func TestSendNotificationEmptyUserID(t *testing.T) {
	payload := NotificationRequest{
		Message: "Test notification",
		UserID:  "",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", status)
	}
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

// TestHealthCheckContentType tests health check returns correct content type
func TestHealthCheckContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestSendNotificationContentType tests notification returns correct content type
func TestSendNotificationContentType(t *testing.T) {
	payload := NotificationRequest{
		Message: "Test notification",
		UserID:  "user123",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestSendNotificationLongMessage tests with long message
func TestSendNotificationLongMessage(t *testing.T) {
	longMessage := ""
	for i := 0; i < 1000; i++ {
		longMessage += "a"
	}

	payload := NotificationRequest{
		Message: longMessage,
		UserID:  "user123",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}
}

// TestSendNotificationSpecialCharacters tests with special characters
func TestSendNotificationSpecialCharacters(t *testing.T) {
	payload := NotificationRequest{
		Message: "Test @#$%^&*()_+-=[]{}|;:,.<>?",
		UserID:  "user@123#456",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}
}

// TestSendNotificationUnicodeCharacters tests with unicode characters
func TestSendNotificationUnicodeCharacters(t *testing.T) {
	payload := NotificationRequest{
		Message: "Hello 世界 🌍 Привет",
		UserID:  "user_unicode_123",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}
}

// BenchmarkSendNotification benchmarks the SendNotification handler
func BenchmarkSendNotification(b *testing.B) {
	payload := NotificationRequest{
		Message: "Benchmark notification",
		UserID:  "user123",
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(SendNotification)
		handler.ServeHTTP(rr, req)
	}
}

// BenchmarkHealthCheck benchmarks the HealthCheck handler
func BenchmarkHealthCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HealthCheck)
		handler.ServeHTTP(rr, req)
	}
}
