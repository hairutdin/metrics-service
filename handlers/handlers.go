package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/storage"
)

type MetricsHandler struct {
	storage storage.MetricsStorage
}

func NewMetricsHandler(s storage.MetricsStorage) *MetricsHandler {
	return &MetricsHandler{storage: s}
}

// HandleGetValue handles GET requests to fetch a specific metric by type and name.
func (h *MetricsHandler) HandleGetValue(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	value, err := h.storage.GetMetric(metricType, metricName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if metricType == "gauge" {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s: %.1f\n", metricName, floatValue)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s: %s\n", metricName, value)))
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

// HandleUpdate handles POST requests to update metrics (gauge/counter).
func (h *MetricsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
	if len(parts) != 3 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	metricType := parts[0]
	metricName := parts[1]
	metricValue := parts[2]

	if metricName == "" {
		http.Error(w, "Metric name is required", http.StatusNotFound)
		return
	}

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(metricName, value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(metricName, value)
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
