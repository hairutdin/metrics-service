package storage

import "github.com/hairutdin/metrics-service/models"

type MetricsStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	UpdateMetricsBatch(metrics []models.Metrics) error
	GetMetric(metricType string, name string) (string, error)
	GetAllMetrics() map[string]string
}
