package main

import (
	"log"
	"os"

	_ "github.com/LucasM4r/repomind/docs"
	"github.com/LucasM4r/repomind/internal/api"
	"github.com/LucasM4r/repomind/internal/api/router"
	"github.com/LucasM4r/repomind/internal/app"
	"github.com/LucasM4r/repomind/internal/config"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("[FATAL] Configuration error: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize application: %v", err)
	}
	defer application.Close()

	address := getEnv("APP_ADDRESS", "0.0.0.0")
	port := getEnv("APP_PORT", "3001")

	r := router.NewRouter(application.RAGService())
	server := api.NewServer(r, address, port)

	if err := server.Start(); err != nil {
		log.Fatalf("[FATAL] Failed to start HTTP server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
