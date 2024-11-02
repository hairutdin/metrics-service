package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/internal/middleware"
	"github.com/hairutdin/metrics-service/storage"
	"github.com/stretchr/testify/assert"
)

func testSetupRouter() *chi.Mux {
	memStorage := storage.NewMemStorage()
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	r := chi.NewRouter()
	r.Use(middleware.GzipCompress)
	r.Use(middleware.GzipDecompress)

	r.Post("/update/", metricsHandler.HandleUpdateJSON)
	r.Post("/value/", metricsHandler.HandleGetValueJSON)
	r.Get("/", metricsHandler.HandleListMetrics)
	r.Get("/ping", handlers.PingHandler(func() error { return nil })) // Mock PingDB to return nil (no error)

	return r
}

func TestUpdateMetric(t *testing.T) {
	router := testSetupRouter()
	metric := handlers.Metrics{
		ID:    "test_metric",
		MType: "gauge",
		Value: func(v float64) *float64 { return &v }(12.5),
	}

	body, err := json.Marshal(metric)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/update/", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")
}

func TestGetValueMetric(t *testing.T) {
	storage := storage.NewMemStorage()
	storage.UpdateGauge("test_metric", 12.5)
	router := setupRouter(storage)

	metric := handlers.Metrics{
		ID:    "test_metric",
		MType: "gauge",
	}

	body, err := json.Marshal(metric)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/value/", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")

	var response handlers.Metrics
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.Value, "Expected a non-nil Value field")
	assert.Equal(t, 12.5, *response.Value, "Expected gauge value 12.5")
}

func TestListMetrics(t *testing.T) {
	storage := storage.NewMemStorage()
	storage.UpdateGauge("gauge_metric", 10.5)
	storage.UpdateCounter("counter_metric", 5)
	router := setupRouter(storage)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")
	assert.Contains(t, rr.Body.String(), "gauge_metric: gauge: 10.500000")
	assert.Contains(t, rr.Body.String(), "counter_metric: counter: 5")
}

func TestPingHandler(t *testing.T) {
	router := testSetupRouter()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")
	assert.Equal(t, "Database connection is OK", rr.Body.String())
}

func TestGzipCompression(t *testing.T) {
	server := httptest.NewServer(middleware.GzipCompress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "compressed response"}`))
	})))
	defer server.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	gzReader, err := gzip.NewReader(resp.Body)
	assert.NoError(t, err)
	defer gzReader.Close()

	body, err := io.ReadAll(gzReader)
	assert.NoError(t, err)

	expectedBody := `{"message": "compressed response"}`
	assert.Equal(t, expectedBody, string(body))
}
