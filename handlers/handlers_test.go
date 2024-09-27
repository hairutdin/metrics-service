package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

func TestHandleGetValue(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(memStorage)

	memStorage.UpdateGauge("test_metric", 10.5)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", metricsHandler.HandleGetValue)

	req, err := http.NewRequest("GET", "/value/gauge/test_metric", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	expected := "test_metric: 10.5\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}

	req, err = http.NewRequest("GET", "/value/gauge/non_existent_metric", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	rr = httptest.NewRecorder()
	metricsHandler.HandleGetValue(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
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
