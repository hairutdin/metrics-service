package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	metrics := &Metrics{}
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

		expectedPath := "/update/gauge/Alloc/100.0"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	sendMetric := func(metricType, metricName string, value float64) {
		url := server.URL + "/update/" + metricType + "/" + metricName + "/" + strconv.FormatFloat(value, 'f', 1, 64)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "text/plain")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	}

	sendMetric("gauge", "Alloc", 100.0)
}
