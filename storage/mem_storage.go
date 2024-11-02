package storage

import (
	"encoding/json"
	"fmt"
	"github.com/hairutdin/metrics-service/models"
	"os"
	"sync"
	"time"
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

func (ms *MemStorage) UpdateMetricsBatch(metrics []models.Metrics) error {
	ms.Lock()
	defer ms.Unlock()
	for _, metric := range metrics {
		if metric.MType == "gauge" {
			ms.Gauges[metric.ID] = *metric.Value
		} else if metric.MType == "counter" {
			ms.Counters[metric.ID] += *metric.Delta
		}
	}
	return nil
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

func (s *MemStorage) SaveMetricsToFile(filePath string) error {
	s.RLock()
	defer s.RUnlock()

	data := struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]int64   `json:"counters"`
	}{
		Gauges:   s.Gauges,
		Counters: s.Counters,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding metrics to file: %v", err)
	}

	return nil
}

func (s *MemStorage) RestoreMetricsFromFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	data := struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]int64   `json:"counters"`
	}{}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("error decoding metrics from file: %v", err)
	}

	s.Gauges = data.Gauges
	s.Counters = data.Counters

	return nil
}

func (s *MemStorage) EnableSyncSaving(filePath string) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			s.SaveMetricsToFile(filePath)
		}
	}()
}
