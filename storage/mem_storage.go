package storage

import (
	"fmt"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

var _ MetricsStorage = (*MemStorage)(nil)

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.Lock()
	defer s.Unlock()
	s.Gauges[name] = value
}

func (s *MemStorage) UpdateCounter(name string, value int64) {
	s.Lock()
	defer s.Unlock()
	s.Counters[name] += value
}

func (s *MemStorage) GetMetric(metricType string, name string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	switch metricType {
	case "gauge":
		if value, exists := s.Gauges[name]; exists {
			return fmt.Sprintf("%f", value), nil
		}
	case "counter":
		if value, exists := s.Counters[name]; exists {
			return fmt.Sprintf("%d", value), nil
		}
	}
	return "", fmt.Errorf("metric not found")
}

func (s *MemStorage) GetAllMetrics() map[string]string {
	s.RLock()
	defer s.RUnlock()

	metrics := make(map[string]string)
	for name, value := range s.Gauges {
		metrics[name] = fmt.Sprintf("gauge: %f", value)
	}
	for name, value := range s.Counters {
		metrics[name] = fmt.Sprintf("counter: %d", value)
	}
	return metrics
}
