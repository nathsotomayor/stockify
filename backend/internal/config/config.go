package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	StockAPIToken string
	ServerPort    string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("ERROR: DATABASE_URL environment variable not set.")
	}

	apiToken := os.Getenv("STOCK_API_TOKEN")
	if apiToken == "" {
		log.Fatal("ERROR: STOCK_API_TOKEN environment variable not set.")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL:   dbURL,
		StockAPIToken: apiToken,
		ServerPort:    ":" + port,
	}
}
