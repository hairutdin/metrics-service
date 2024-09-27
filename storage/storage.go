package storage

type MetricsStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetMetric(metricType string, name string) (string, error)
	GetAllMetrics() map[string]string
}
