package storage

import (
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
