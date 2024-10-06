package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

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
	metric := Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: delta,
		Value: value,
	}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	url := "http://" + serverAddress + "/update/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
			sendMetric("gauge", "Alloc", nil, &metrics.Alloc, serverAddress)
			sendMetric("gauge", "BuckHashSys", nil, &metrics.BuckHashSys, serverAddress)
			sendMetric("gauge", "Frees", nil, &metrics.Frees, serverAddress)
			sendMetric("gauge", "GCCPUFraction", nil, &metrics.GCCPUFraction, serverAddress)
			sendMetric("gauge", "GCSys", nil, &metrics.GCSys, serverAddress)
			sendMetric("gauge", "HeapAlloc", nil, &metrics.HeapAlloc, serverAddress)
			sendMetric("gauge", "HeapIdle", nil, &metrics.HeapIdle, serverAddress)
			sendMetric("gauge", "HeapInuse", nil, &metrics.HeapInuse, serverAddress)
			sendMetric("gauge", "HeapObjects", nil, &metrics.HeapObjects, serverAddress)
			sendMetric("gauge", "HeapReleased", nil, &metrics.HeapReleased, serverAddress)
			sendMetric("gauge", "HeapSys", nil, &metrics.HeapSys, serverAddress)
			sendMetric("gauge", "LastGC", nil, &metrics.LastGC, serverAddress)
			sendMetric("gauge", "Lookups", nil, &metrics.Lookups, serverAddress)
			sendMetric("gauge", "MCacheInuse", nil, &metrics.MCacheInuse, serverAddress)
			sendMetric("gauge", "MCacheSys", nil, &metrics.MCacheSys, serverAddress)
			sendMetric("gauge", "MSpanInuse", nil, &metrics.MSpanInuse, serverAddress)
			sendMetric("gauge", "MSpanSys", nil, &metrics.MSpanSys, serverAddress)
			sendMetric("gauge", "Mallocs", nil, &metrics.Mallocs, serverAddress)
			sendMetric("gauge", "NextGC", nil, &metrics.NextGC, serverAddress)
			sendMetric("gauge", "NumForcedGC", nil, &metrics.NumForcedGC, serverAddress)
			sendMetric("gauge", "NumGC", nil, &metrics.NumGC, serverAddress)
			sendMetric("gauge", "OtherSys", nil, &metrics.OtherSys, serverAddress)
			sendMetric("gauge", "PauseTotalNs", nil, &metrics.PauseTotalNs, serverAddress)
			sendMetric("gauge", "StackInuse", nil, &metrics.StackInuse, serverAddress)
			sendMetric("gauge", "StackSys", nil, &metrics.StackSys, serverAddress)
			sendMetric("gauge", "Sys", nil, &metrics.Sys, serverAddress)
			sendMetric("gauge", "TotalAlloc", nil, &metrics.TotalAlloc, serverAddress)

			sendMetric("counter", "PollCount", &metrics.PollCount, nil, serverAddress)
			sendMetric("gauge", "RandomValue", nil, &metrics.RandomValue, serverAddress)
			time.Sleep(time.Duration(reportInterval) * time.Second)
		}
	}()

	select {}
}
