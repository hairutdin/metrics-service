package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/storage"
	"net/http"
	"os"
)

func main() {
	flagServerAddress := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	envServerAddress := os.Getenv("SERVER_ADDRESS")

	serverAddress := *flagServerAddress
	if envServerAddress != "" {
		serverAddress = envServerAddress
	}

	if len(flag.Args()) > 0 {
		fmt.Printf("Error: Unknown flags or arguments: %v\n", flag.Args())
		os.Exit(1)
	}

	memStorage := storage.NewMemStorage()
	r := chi.NewRouter()

	metricsHandler := handlers.NewMetricsHandler(memStorage)

	r.Post("/update/", metricsHandler.HandleUpdate)
	r.Get("/value/{type}/{name}", metricsHandler.HandleGetValue)
	r.Get("/", metricsHandler.HandleListMetrics)

	fmt.Printf("Server is running at http://%s\n", serverAddress)
	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
