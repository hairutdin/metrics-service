package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/internal/middleware"
	"github.com/hairutdin/metrics-service/storage"
	"github.com/sirupsen/logrus"
)

func startMetricSaver(interval int, filePath string, storage *storage.MemStorage) {
	if interval == 0 {
		storage.EnableSyncSaving(filePath)
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			<-ticker.C
			storage.SaveMetricsToFile(filePath)
		}
	}()

}

func main() {
	flagServerAddress := flag.String("a", "localhost:8080", "server address")
	flagStoreInterval := flag.Int("i", 300, "Store interval in seconds")
	flagFilePath := flag.String("f", "/tmp/metrics-db.json", "File path to save metrics")
	flagRestore := flag.Bool("r", true, "Restore metrics from file on startup")
	flag.Parse()

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	envStoreInterval := os.Getenv("STORE_INTERVAL")
	envFilePath := os.Getenv("FILE_STORAGE_PATH")
	envRestore := os.Getenv("RESTORE")

	serverAddress := *flagServerAddress
	if envServerAddress != "" {
		serverAddress = envServerAddress
	}
	storeInterval := *flagStoreInterval
	if envStoreInterval != "" {
		interval, err := strconv.Atoi(envStoreInterval)
		if err == nil {
			storeInterval = interval
		}
	}
	filePath := *flagFilePath
	if envFilePath != "" {
		filePath = envFilePath
	}
	restore := *flagRestore
	if envRestore != "" {
		restoreBool, err := strconv.ParseBool(envRestore)
		if err == nil {
			restore = restoreBool
		}
	}

	if len(flag.Args()) > 0 {
		fmt.Printf("Error: Unknown flags or arguments: %v\n", flag.Args())
		os.Exit(1)
	}

	memStorage := storage.NewMemStorage()

	if restore {
		err := memStorage.RestoreMetricsFromFile(filePath)
		if err != nil {
			fmt.Printf("Error restoring metrics from file: %v\n", err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	startMetricSaver(storeInterval, filePath, memStorage)

	r := chi.NewRouter()
	logger := logrus.New()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.GzipDecompress)
	r.Use(middleware.GzipCompress)

	metricsHandler := handlers.NewMetricsHandler(memStorage)
	r.Post("/update/", metricsHandler.HandleUpdateJSON)
	r.Post("/value/", metricsHandler.HandleGetValueJSON)
	r.Get("/", metricsHandler.HandleListMetrics)

	go func() {
		fmt.Printf("Server is running at http://%s\n", serverAddress)
		err := http.ListenAndServe(serverAddress, r)
		if err != nil {
			fmt.Printf("Server failed to start: %v\n", err)
		}
	}()

	<-stop

	fmt.Println("Shutting down server... Saving metrics.")
	memStorage.SaveMetricsToFile(filePath)
	os.Exit(0)
}
