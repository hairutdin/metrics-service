package main

import (
	"bytes"
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

func (m *Metrics) collectMetrics() {
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

func sendMetric(metricType, metricName string, value float64, serverAddress string) {
	url := "http://" + serverAddress + "/update/" + metricType + "/" + metricName + "/" + strconv.FormatFloat(value, 'f', 1, 64)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func main() {
	reportInterval := flag.Int("r", 10, "Report interval in seconds")
	pollInterval := flag.Int("p", 2, "Poll interval in seconds")
	serverAddress := flag.String("a", "localhost:8080", "server address")

	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Error: Unknown flags or arguments: %v\n", flag.Args())
		os.Exit(1)
	}

	metrics := &Metrics{}

	go func() {
		for {
			metrics.collectMetrics()
			time.Sleep(time.Duration(*pollInterval) * time.Second)
		}
	}()

	go func() {
		for {
			sendMetric("gauge", "Alloc", metrics.Alloc, *serverAddress)
			sendMetric("gauge", "BuckHashSys", metrics.BuckHashSys, *serverAddress)
			sendMetric("gauge", "Frees", metrics.Frees, *serverAddress)
			sendMetric("gauge", "GCCPUFraction", metrics.GCCPUFraction, *serverAddress)
			sendMetric("gauge", "GCSys", metrics.GCSys, *serverAddress)
			sendMetric("gauge", "HeapAlloc", metrics.HeapAlloc, *serverAddress)
			sendMetric("gauge", "HeapIdle", metrics.HeapIdle, *serverAddress)
			sendMetric("gauge", "HeapInuse", metrics.HeapInuse, *serverAddress)
			sendMetric("gauge", "HeapObjects", metrics.HeapObjects, *serverAddress)
			sendMetric("gauge", "HeapReleased", metrics.HeapReleased, *serverAddress)
			sendMetric("gauge", "HeapSys", metrics.HeapSys, *serverAddress)
			sendMetric("gauge", "LastGC", metrics.LastGC, *serverAddress)
			sendMetric("gauge", "Lookups", metrics.Lookups, *serverAddress)
			sendMetric("gauge", "MCacheInuse", metrics.MCacheInuse, *serverAddress)
			sendMetric("gauge", "MCacheSys", metrics.MCacheSys, *serverAddress)
			sendMetric("gauge", "MSpanInuse", metrics.MSpanInuse, *serverAddress)
			sendMetric("gauge", "MSpanSys", metrics.MSpanSys, *serverAddress)
			sendMetric("gauge", "Mallocs", metrics.Mallocs, *serverAddress)
			sendMetric("gauge", "NextGC", metrics.NextGC, *serverAddress)
			sendMetric("gauge", "NumForcedGC", metrics.NumForcedGC, *serverAddress)
			sendMetric("gauge", "NumGC", metrics.NumGC, *serverAddress)
			sendMetric("gauge", "OtherSys", metrics.OtherSys, *serverAddress)
			sendMetric("gauge", "PauseTotalNs", metrics.PauseTotalNs, *serverAddress)
			sendMetric("gauge", "StackInuse", metrics.StackInuse, *serverAddress)
			sendMetric("gauge", "StackSys", metrics.StackSys, *serverAddress)
			sendMetric("gauge", "Sys", metrics.Sys, *serverAddress)
			sendMetric("gauge", "TotalAlloc", metrics.TotalAlloc, *serverAddress)

			sendMetric("counter", "PollCount", float64(metrics.PollCount), *serverAddress)
			sendMetric("gauge", "RandomValue", metrics.RandomValue, *serverAddress)
			time.Sleep(time.Duration(*reportInterval) * time.Second)
		}
	}()

	select {}
}
