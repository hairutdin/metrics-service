package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/storage"
	"net/http"
)

func main() {
	memStorage := storage.NewMemStorage()
	r := chi.NewRouter()

	metricsHandler := handlers.NewMetricsHandler(memStorage)

	r.Post("/update/", metricsHandler.HandleUpdate)
	r.Get("/value/{type}/{name}", metricsHandler.HandleGetValue)
	r.Get("/", metricsHandler.HandleListMetrics)

	fmt.Println("Server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
