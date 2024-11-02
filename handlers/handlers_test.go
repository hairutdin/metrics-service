package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/models"
	"github.com/hairutdin/metrics-service/storage"
	"github.com/stretchr/testify/assert"
)

func TestPingHandlerSuccess(t *testing.T) {
	mockPingDB := func() error {
		return nil
	}

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr := httptest.NewRecorder()

	handler := PingHandler(mockPingDB)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	expected := "Database connection is OK"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}
}

func TestPingHandlerFailure(t *testing.T) {
	mockPingDB := func() error {
		return errors.New("database connection failed")
	}

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr := httptest.NewRecorder()

	handler := PingHandler(mockPingDB)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
	}

	expected := "Database connection failed\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}
}

func TestHandleUpdateJSON(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(memStorage)

	metric := models.Metrics{
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

	metric := models.Metrics{
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
	expected := models.Metrics{
		Value: new(float64),
	}
	*expected.Value = 10.5

	var response models.Metrics
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

func TestUpdateMetricsBatchHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	handler := NewMetricsHandler(memStorage)

	r := chi.NewRouter()
	r.Post("/updates/", handler.HandleBatchUpdate)

	metrics := []models.Metrics{
		{ID: "batch_gauge", MType: "gauge", Value: func(v float64) *float64 { val := v; return &val }(52.5)},
		{ID: "batch_counter", MType: "counter", Delta: func(d int64) *int64 { delta := d; return &delta }(20)},
	}

	jsonData, err := json.Marshal(metrics)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")

	gaugeValue, _ := memStorage.GetMetric("gauge", "batch_gauge")
	counterValue, _ := memStorage.GetMetric("counter", "batch_counter")
	assert.Equal(t, "52.500000", gaugeValue)
	assert.Equal(t, "20", counterValue)
}
