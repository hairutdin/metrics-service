package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hairutdin/metrics-service/storage"
)

func TestHandleUpdate(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(memStorage)

	req, err := http.NewRequest("POST", "/update/gauge/test_metric/10.5", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr := httptest.NewRecorder()
	metricsHandler.HandleUpdate(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	if val, ok := memStorage.Gauges["test_metric"]; !ok || val != 10.5 {
		t.Errorf("Expected gauge value 10.5, got %v", val)
	}

	req, err = http.NewRequest("POST", "/update/counter/test_counter/10.5", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr = httptest.NewRecorder()
	metricsHandler.HandleUpdate(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid metric type, got %d", status)
	}
}
