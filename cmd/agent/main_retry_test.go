package main

import (
	"errors"
	"github.com/hairutdin/metrics-service/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestSendMetricBatch_RetryOnNetworkError(t *testing.T) {
	mockTransport := &MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("temporary network error")
		},
	}

	client := &http.Client{Transport: mockTransport}

	metrics := []models.Metrics{
		{ID: "Alloc", MType: "gauge", Value: func(v float64) *float64 { return &v }(42.0)},
	}

	err := sendMetricsBatch(metrics, "localhost:8080", client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "temporary network error")
}
