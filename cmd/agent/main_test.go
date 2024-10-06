package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	metrics := &RuntimeMetrics{}
	metrics.collectMetrics()

	if metrics.Alloc < 0 {
		t.Errorf("Alloc should not be less than 0")
	}
}

func TestSendMetric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content type, got %s", r.Header.Get("Content-Type"))
		}

		var metric Metrics
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			t.Fatalf("Failed to decode JSON: %v", err)
		}

		if metric.ID != "Alloc" {
			t.Errorf("Expected metric ID 'Alloc', got %s", metric.ID)
		}
		if metric.MType != "gauge" {
			t.Errorf("Expected metric type 'gauge', got %s", metric.MType)
		}
		if metric.Value == nil || *metric.Value != 100.0 {
			t.Errorf("Expected metric value 100.0, got %v", metric.Value)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	sendMetric := func(metricType, metricName string, value float64) {
		metric := Metrics{
			ID:    metricName,
			MType: metricType,
			Value: &value,
		}
		data, err := json.Marshal(metric)
		if err != nil {
			t.Fatalf("Failed to marshal metric: %v", err)
		}

		url := server.URL + "/update/"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	}

	sendMetric("gauge", "Alloc", 100.0)
}
