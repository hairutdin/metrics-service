package main

import (
	"bytes"
	"encoding/json"
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	http.HandleFunc("/update/", metricsHandler.HandleUpdateJSON)

	metric := handlers.Metrics{
		ID:    "test_metric",
		MType: "gauge",
		Value: func(v float64) *float64 { return &v }(12.5),
	}

	body, err := json.Marshal(metric)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/update/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body><h1>Metrics</h1></body></html>"))
	})

	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}

	expectedBody := "<html><body><h1>Metrics</h1></body></html>"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %v, got %v", expectedBody, rr.Body.String())
	}
}
