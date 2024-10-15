package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hairutdin/metrics-service/storage"
)

type Metrics struct {
	ID    string   `json:"id"`              // name of the metric
	MType string   `json:"type"`            // gauge or counter
	Delta *int64   `json:"delta,omitempty"` // value for counter
	Value *float64 `json:"value,omitempty"` // value for gauge
}

type MetricsHandler struct {
	storage storage.MetricsStorage
}

func NewMetricsHandler(s storage.MetricsStorage) *MetricsHandler {
	return &MetricsHandler{storage: s}
}

// HandleUpdateJSON handles POST requests to update metrics in JSON format
func (h *MetricsHandler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var metric Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			http.Error(w, "Gauge value is required", http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(metric.ID, *metric.Value)
	case "counter":
		if metric.Delta == nil {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}
}

func (h *MetricsHandler) HandleGetValueJSON(w http.ResponseWriter, r *http.Request) {
	var metric Metrics
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var value string
	switch metric.MType {
	case "gauge":
		value, err = h.storage.GetMetric("gauge", metric.ID)
		if err != nil {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		val, _ := strconv.ParseFloat(value, 64)
		metric.Value = &val
	case "counter":
		value, err = h.storage.GetMetric("counter", metric.ID)
		if err != nil {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		delta, _ := strconv.ParseInt(value, 10, 64)
		metric.Delta = &delta
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(metric)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// HandleListMetrics handles GET requests to list all known metrics in HTML format.
func (h *MetricsHandler) HandleListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.storage.GetAllMetrics()

	w.Header().Set("Content-Type", "text/html")
	html := "<html><body><h1>Metrics</h1><ul>"
	for name, value := range metrics {
		html += fmt.Sprintf("<li>%s: %v</li>", name, value)
	}
	html += "</ul></body></html>"

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

type PingDBFunc func() error

func PingHandler(pingDB PingDBFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pingDB()
		if err != nil {
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection is OK"))
	}
}
