package storage

import (
	"testing"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()

	storage.UpdateGauge("testGauge", 10.5)
	if val, ok := storage.Gauges["testGauge"]; !ok || val != 10.5 {
		t.Errorf("Expected gauge value 10.5, got %v", val)
	}

	storage.UpdateCounter("testCounter", 5)
	if val, ok := storage.Counters["testCounter"]; !ok || val != 5 {
		t.Errorf("Expected counter value 5, got %v", val)
	}

	storage.UpdateCounter("testCounter", 3)
	if val, ok := storage.Counters["testCounter"]; !ok || val != 8 {
		t.Errorf("Expected counter value 8, got %v", val)
	}

	value, err := storage.GetMetric("gauge", "testGauge")
	if err != nil || value != "10.500000" {
		t.Errorf("Expected gauge value 10.500000, got %s", value)
	}

	value, err = storage.GetMetric("counter", "testCounter")
	if err != nil || value != "8" {
		t.Errorf("Expected counter value 8, got %s", value)
	}

	_, err = storage.GetMetric("gauge", "nonExistentGauge")
	if err == nil {
		t.Errorf("Expected error for non-existent gauge metric, got none")
	}

	storage.UpdateGauge("gauge2", 99.9)
	storage.UpdateCounter("counter2", 123)
	metrics := storage.GetAllMetrics()

	expectedMetrics := map[string]string{
		"testGauge":   "gauge: 10.500000",
		"gauge2":      "gauge: 99.900000",
		"testCounter": "counter: 8",
		"counter2":    "counter: 123",
	}

	for name, expected := range expectedMetrics {
		if val, ok := metrics[name]; !ok || val != expected {
			t.Errorf("Expected metric %s: '%s', got '%s'", name, expected, val)
		}
	}
}
