package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/storage"
)

func TestHandleUpdateJSON(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(memStorage)

	metric := Metrics{
		ID:    "test_metric",
		MType: "gauge",
		Value: new(float64),
	}
	*metric.Value = 10.5

	body, _ := json.Marshal(metric)
	req, err := http.NewRequest("POST", "/update/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	metricsHandler.HandleUpdateJSON(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	if val, ok := memStorage.Gauges["test_metric"]; !ok || val != 10.5 {
		t.Errorf("Expected gauge value 10.5, got %v", val)
	}
}

func TestHandleGetValue(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(memStorage)

	memStorage.UpdateGauge("test_metric", 10.5)

	r := chi.NewRouter()
	r.Post("/value/", metricsHandler.HandleGetValueJSON)

	metric := Metrics{
		ID:    "test_metric",
		MType: "gauge",
	}
	body, _ := json.Marshal(metric)

	req, err := http.NewRequest("POST", "/value/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}
	expected := Metrics{
		Value: new(float64),
	}
	*expected.Value = 10.5

	var response Metrics
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	if *response.Value != *expected.Value {
		t.Errorf("Expected value %v, got %v", *expected.Value, *response.Value)
	}
}

func TestHandleListMetrics(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(memStorage)

	memStorage.UpdateGauge("test_gauge", 5.0)
	memStorage.UpdateCounter("test_counter", 15)

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr := httptest.NewRecorder()
	metricsHandler.HandleListMetrics(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	expected := "<html><body><h1>Metrics</h1><ul><li>test_gauge: gauge: 5.000000</li><li>test_counter: counter: 15</li></ul></body></html>"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}
}
