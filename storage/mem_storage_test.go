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
}
