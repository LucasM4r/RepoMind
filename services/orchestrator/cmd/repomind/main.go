package main

import (
	"log"

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

	application.Run()
}
