package main

import (
	"fmt"
	"github.com/hairutdin/metrics-service/handlers"
	"github.com/hairutdin/metrics-service/storage"
	"net/http"
)

func main() {
	memStorage := storage.NewMemStorage()

	metricsHandler := handlers.NewMetricsHandler(memStorage)

	http.HandleFunc("/update/", metricsHandler.HandleUpdate)

	fmt.Println("Server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
