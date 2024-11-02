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
	"github.com/hairutdin/metrics-service/internal/db"
	"github.com/hairutdin/metrics-service/internal/middleware"
	"github.com/hairutdin/metrics-service/storage"
	metricsStorage "github.com/hairutdin/metrics-service/storage"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func initializeStorage(dsn, filePath string, restore bool) storage.MetricsStorage {
	if dsn != "" {
		conn, err := db.ConnectToDB(dsn)
		if err == nil {
			fmt.Println("Using PostgreSQL storage.")
			db.InitializeTables(conn)
			return storage.NewPostgresStorage(conn)
		}
		fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
		defer db.CloseDB(conn)
	}

	memStorage := storage.NewMemStorage()
	if filePath != "" && restore {
		err := memStorage.RestoreMetricsFromFile(filePath)
		if err != nil {
			fmt.Printf("Error restoring metrics from file: %v\n", err)
		} else {
			fmt.Println("Using file storage.")
			return memStorage
		}
	}
	fmt.Println("Using in-memory storage.")
	return memStorage
}

func setupRouter(storage storage.MetricsStorage) *chi.Mux {
	metricsHandler := handlers.NewMetricsHandler(storage)

	r := chi.NewRouter()
	logger := logrus.New()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.GzipDecompress)
	r.Use(middleware.GzipCompress)

	r.Post("/update/", metricsHandler.HandleUpdateJSON)
	r.Post("/updates/", metricsHandler.HandleBatchUpdate)
	r.Post("/value/", metricsHandler.HandleGetValueJSON)
	r.Get("/", metricsHandler.HandleListMetrics)
	r.Get("/ping", handlers.PingHandler(func() error {
		if pgStorage, ok := storage.(*metricsStorage.PostgresStorage); ok {
			return db.PingDB(pgStorage.DB)
		}
		return nil
	}))

	return r
}

func startMetricSaver(interval int, filePath string, storage storage.MetricsStorage) {
	if interval == 0 {
		if memStorage, ok := storage.(*metricsStorage.MemStorage); ok {
			memStorage.EnableSyncSaving(filePath)
		}
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			<-ticker.C
			if memStorage, ok := storage.(*metricsStorage.MemStorage); ok {
				memStorage.SaveMetricsToFile(filePath)
			}
		}
	}()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

func main() {
	flagServerAddress := flag.String("a", "localhost:8080", "server address")
	flagStoreInterval := flag.Int("i", 300, "Store interval in seconds")
	flagFilePath := flag.String("f", "/tmp/metrics-db.json", "File path to save metrics")
	flagRestore := flag.Bool("r", true, "Restore metrics from file on startup")
	flagDSN := flag.String("d", "postgres://postgres:berlin@localhost:5432/testdb?sslmode=disable",
		"PostgreSQL DSN for database connection")
	flag.Parse()

	serverAddress := getEnv("SERVER_ADDRESS", *flagServerAddress)
	storeInterval := getEnvInt("STORE_INTERVAL", *flagStoreInterval)
	filePath := getEnv("FILE_STORAGE_PATH", *flagFilePath)
	restore := getEnvBool("RESTORE", *flagRestore)
	dsn := getEnv("DATABASE_DSN", *flagDSN)

	var metricsStorage storage.MetricsStorage
	var conn *pgx.Conn
	var err error

	if dsn != "" {
		conn, err = db.ConnectToDB(dsn)
		if err == nil {
			metricsStorage = storage.NewPostgresStorage(conn)
			fmt.Println("Using PostgreSQL storage.")
		} else {
			fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
		}
	}

	if metricsStorage == nil && filePath != "" {
		memStorage := storage.NewMemStorage()
		if restore {
			if err := memStorage.RestoreMetricsFromFile(filePath); err != nil {
				fmt.Printf("Error restoring metrics from file: %v\n", err)
			}
		}
		metricsStorage = memStorage
		fmt.Println("Using file-based storage.")
	}

	if metricsStorage == nil {
		metricsStorage = storage.NewMemStorage()
		fmt.Println("Using in-memory storage.")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	startMetricSaver(storeInterval, filePath, metricsStorage)

	r := setupRouter(metricsStorage)

	go func() {
		fmt.Printf("Server is running at http://%s\n", serverAddress)
		err := http.ListenAndServe(serverAddress, r)
		if err != nil {
			fmt.Printf("Server failed to start: %v\n", err)
		}
	}()

	<-stop

	fmt.Println("Shutting down server... Saving metrics.")
	if memStorage, ok := metricsStorage.(*storage.MemStorage); ok {
		memStorage.SaveMetricsToFile(filePath)
	}
	if conn != nil {
		db.CloseDB(conn)
	}
}
