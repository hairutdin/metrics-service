package main

import (
	"bytes"
	"math/rand"
	"net/http"
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

func sendMetric(metricType, metricName string, value float64) {
	url := "http://localhost:8080/update/" + metricType + "/" + metricName + "/" + strconv.FormatFloat(value, 'f', 1, 64)
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
	metrics := &Metrics{}

	go func() {
		for {
			metrics.collectMetrics()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			sendMetric("gauge", "Alloc", metrics.Alloc)
			sendMetric("gauge", "BuckHashSys", metrics.BuckHashSys)
			sendMetric("gauge", "Frees", metrics.Frees)
			sendMetric("gauge", "GCCPUFraction", metrics.GCCPUFraction)
			sendMetric("gauge", "GCSys", metrics.GCSys)
			sendMetric("gauge", "HeapAlloc", metrics.HeapAlloc)
			sendMetric("gauge", "HeapIdle", metrics.HeapIdle)
			sendMetric("gauge", "HeapInuse", metrics.HeapInuse)
			sendMetric("gauge", "HeapObjects", metrics.HeapObjects)
			sendMetric("gauge", "HeapReleased", metrics.HeapReleased)
			sendMetric("gauge", "HeapSys", metrics.HeapSys)
			sendMetric("gauge", "LastGC", metrics.LastGC)
			sendMetric("gauge", "Lookups", metrics.Lookups)
			sendMetric("gauge", "MCacheInuse", metrics.MCacheInuse)
			sendMetric("gauge", "MCacheSys", metrics.MCacheSys)
			sendMetric("gauge", "MSpanInuse", metrics.MSpanInuse)
			sendMetric("gauge", "MSpanSys", metrics.MSpanSys)
			sendMetric("gauge", "Mallocs", metrics.Mallocs)
			sendMetric("gauge", "NextGC", metrics.NextGC)
			sendMetric("gauge", "NumForcedGC", metrics.NumForcedGC)
			sendMetric("gauge", "NumGC", metrics.NumGC)
			sendMetric("gauge", "OtherSys", metrics.OtherSys)
			sendMetric("gauge", "PauseTotalNs", metrics.PauseTotalNs)
			sendMetric("gauge", "StackInuse", metrics.StackInuse)
			sendMetric("gauge", "StackSys", metrics.StackSys)
			sendMetric("gauge", "Sys", metrics.Sys)
			sendMetric("gauge", "TotalAlloc", metrics.TotalAlloc)

			sendMetric("counter", "PollCount", float64(metrics.PollCount))
			sendMetric("gauge", "RandomValue", metrics.RandomValue)
			time.Sleep(10 * time.Second)
		}
	}()

	select {}
}
