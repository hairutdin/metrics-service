package main

import (
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/storage"
	"net/http"
	"net/http/httptest"
	"strings"
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

	expectedBody := "<html>"
	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, rr.Body.String())
	}
}
