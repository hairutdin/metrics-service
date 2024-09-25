package main

import (
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	http.HandleFunc("/update/", metricsHandler.HandleUpdate)

	req, err := http.NewRequest("POST", "/update/gauge/test_metric/12.5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}
}
