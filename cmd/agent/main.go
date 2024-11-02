package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/hairutdin/metrics-service/models"
)

// type Metrics struct {
// 	ID    string   `json:"id"`
// 	MType string   `json:"type"`
// 	Delta *int64   `json:"delta,omitempty"`
// 	Value *float64 `json:"value,omitempty"`
// }

type RuntimeMetrics struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64

	PollCount   int64
	RandomValue float64
}

func randomValue() float64 {
	return rand.Float64()
}

func (m *RuntimeMetrics) collectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.Alloc = float64(memStats.Alloc)
	m.BuckHashSys = float64(memStats.BuckHashSys)
	m.Frees = float64(memStats.Frees)
	m.GCCPUFraction = float64(memStats.GCCPUFraction)
	m.GCSys = float64(memStats.GCSys)
	m.HeapAlloc = float64(memStats.HeapAlloc)
	m.HeapIdle = float64(memStats.HeapIdle)
	m.HeapInuse = float64(memStats.HeapInuse)
	m.HeapObjects = float64(memStats.HeapObjects)
	m.HeapReleased = float64(memStats.HeapReleased)
	m.HeapSys = float64(memStats.HeapSys)
	m.LastGC = float64(memStats.LastGC)
	m.Lookups = float64(memStats.Lookups)
	m.MCacheInuse = float64(memStats.MCacheInuse)
	m.MCacheSys = float64(memStats.MCacheSys)
	m.MSpanInuse = float64(memStats.MSpanInuse)
	m.MSpanSys = float64(memStats.MSpanSys)
	m.Mallocs = float64(memStats.Mallocs)
	m.NextGC = float64(memStats.NextGC)
	m.NumForcedGC = float64(memStats.NumForcedGC)
	m.NumGC = float64(memStats.NumGC)
	m.OtherSys = float64(memStats.OtherSys)
	m.PauseTotalNs = float64(memStats.PauseTotalNs)
	m.StackInuse = float64(memStats.StackInuse)
	m.StackSys = float64(memStats.StackSys)
	m.Sys = float64(memStats.Sys)
	m.TotalAlloc = float64(memStats.TotalAlloc)

	m.PollCount++
	m.RandomValue = randomValue()
}

func sendMetric(metricType, metricName string, delta *int64, value *float64, serverAddress string) {
	metric := models.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: delta,
		Value: value,
	}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)

	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		panic(fmt.Sprintf("Failed to write to gzip writer: %v", err))
	}

	if err := gzipWriter.Close(); err != nil {
		panic(fmt.Sprintf("Failed to close gzip writer: %v", err))
	}

	url := "http://" + serverAddress + "/update/"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func sendMetricsBatch(metrics []models.Metrics, serverAddress string) {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal metrics: %v", err))
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		panic(fmt.Sprintf("Failed to write to gzip writer: %v", err))
	}

	if err := gzipWriter.Close(); err != nil {
		panic(fmt.Sprintf("Failed to close gzip writer: %v", err))
	}

	url := "http://" + serverAddress + "/updates/"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send metrics batch: %v\n", err)
		return
	}
	defer resp.Body.Close()

}

func main() {
	flagReportInterval := flag.Int("r", 10, "Report interval in seconds")
	flagPollInterval := flag.Int("p", 2, "Poll interval in seconds")
	flagServerAddress := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")
	envServerAddress := os.Getenv("SERVER_ADDRESS")

	reportInterval := *flagReportInterval
	if envReportInterval != "" {
		if value, err := strconv.Atoi(envReportInterval); err == nil {
			reportInterval = value
		} else {
			fmt.Printf("Error: Invalid value for REPORT_INTERVAL: %v\n", envReportInterval)
			os.Exit(1)
		}
	}

	pollInterval := *flagPollInterval
	if envPollInterval != "" {
		if value, err := strconv.Atoi(envPollInterval); err == nil {
			pollInterval = value
		} else {
			fmt.Printf("Error: Invalid value for POLL_INTERVAL: %v\n", envPollInterval)
			os.Exit(1)
		}
	}

	serverAddress := *flagServerAddress
	if envServerAddress != "" {
		serverAddress = envServerAddress
	}

	if len(flag.Args()) > 0 {
		fmt.Printf("Error: Unknown flags or arguments: %v\n", flag.Args())
		os.Exit(1)
	}

	metrics := &RuntimeMetrics{}
	go func() {
		for {
			metrics.collectMetrics()
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			metricList := []models.Metrics{
				{ID: "Alloc", MType: "gauge", Value: &metrics.Alloc},
				{ID: "BuckHashSys", MType: "gauge", Value: &metrics.BuckHashSys},
				{ID: "Frees", MType: "gauge", Value: &metrics.Frees},
				{ID: "GCCPUFraction", MType: "gauge", Value: &metrics.GCCPUFraction},
				{ID: "GCSys", MType: "gauge", Value: &metrics.GCSys},
				{ID: "HeapAlloc", MType: "gauge", Value: &metrics.HeapAlloc},
				{ID: "HeapIdle", MType: "gauge", Value: &metrics.HeapIdle},
				{ID: "HeapInuse", MType: "gauge", Value: &metrics.HeapInuse},
				{ID: "HeapObjects", MType: "gauge", Value: &metrics.HeapObjects},
				{ID: "HeapReleased", MType: "gauge", Value: &metrics.HeapReleased},
				{ID: "HeapSys", MType: "gauge", Value: &metrics.HeapSys},
				{ID: "LastGC", MType: "gauge", Value: &metrics.LastGC},
				{ID: "Lookups", MType: "gauge", Value: &metrics.Lookups},
				{ID: "MCacheInuse", MType: "gauge", Value: &metrics.MCacheInuse},
				{ID: "MCacheSys", MType: "gauge", Value: &metrics.MCacheSys},
				{ID: "MSpanInuse", MType: "gauge", Value: &metrics.MSpanInuse},
				{ID: "MSpanSys", MType: "gauge", Value: &metrics.MSpanSys},
				{ID: "Mallocs", MType: "gauge", Value: &metrics.Mallocs},
				{ID: "NextGC", MType: "gauge", Value: &metrics.NextGC},
				{ID: "NumForcedGC", MType: "gauge", Value: &metrics.NumForcedGC},
				{ID: "NumGC", MType: "gauge", Value: &metrics.NumGC},
				{ID: "OtherSys", MType: "gauge", Value: &metrics.OtherSys},
				{ID: "PauseTotalNs", MType: "gauge", Value: &metrics.PauseTotalNs},
				{ID: "StackInuse", MType: "gauge", Value: &metrics.StackInuse},
				{ID: "StackSys", MType: "gauge", Value: &metrics.StackSys},
				{ID: "Sys", MType: "gauge", Value: &metrics.Sys},
				{ID: "TotalAlloc", MType: "gauge", Value: &metrics.TotalAlloc},
				{ID: "PollCount", MType: "counter", Delta: &metrics.PollCount},
				{ID: "RandomValue", MType: "gauge", Value: &metrics.RandomValue},
			}

			sendMetricsBatch(metricList, serverAddress)
			time.Sleep(time.Duration(reportInterval) * time.Second)
		}
	}()

	select {}
}
