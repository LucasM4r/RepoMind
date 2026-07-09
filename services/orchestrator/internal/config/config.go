package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	GithubToken        string
	MaxWorkers         int
	AIEngineURL        string
	MaxConcurrentFiles int
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[INFO] .env file not found, relying on system environment variables.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		log.Printf("[WARNING] GITHUB_TOKEN is not set. Public repositories will be fetched anonymously.")
	}

	aiURL := os.Getenv("AI_ENGINE_URL")
	if aiURL == "" {
		aiURL = "localhost:50051"
	}

	maxWorkers := 3

	if envWorkers := os.Getenv("MAX_WORKERS"); envWorkers != "" {
		if parsed, err := strconv.Atoi(envWorkers); err == nil && parsed > 0 {
			maxWorkers = parsed
		} else {
			log.Printf("[WARNING] MAX_WORKERS invalid '%s'. Using default: %d\n", envWorkers, maxWorkers)
		}
	}

	return &Config{
		DatabaseURL:        dbURL,
		GithubToken:        ghToken,
		MaxWorkers:         maxWorkers,
		AIEngineURL:        aiURL,
		MaxConcurrentFiles: 5,
	}, nil
}
