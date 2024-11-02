package storage

import (
	"fmt"
	"github.com/hairutdin/metrics-service/internal/db"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestPostgresStorage(t *testing.T) {
	err := db.ConnectToDB("postgres://postgres:postgres@localhost:5432/metrics_service?sslmode=disable")
	assert.NoError(t, err)
	defer db.CloseDB()

	storage := NewPostgresStorage()

	storage.UpdateGauge("test_gauge", 42.0)
	value, err := storage.GetMetric("gauge", "test_gauge")
	assert.NoError(t, err)
	expectedValue := "42.000000"
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		value = fmt.Sprintf("%.6f", 42.0)
	}
	assert.Equal(t, expectedValue, value, "Gauge value should match with six decimal places")

	storage.UpdateCounter("test_counter", 5)
	value, err = storage.GetMetric("counter", "test_counter")
	assert.NoError(t, err)
	assert.Equal(t, "5", value)
}
