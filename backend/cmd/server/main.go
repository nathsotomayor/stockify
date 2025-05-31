package main

import (
	"log"
	"net/http"

	"stockify/internal/api"
	"stockify/internal/config"
	"stockify/internal/database"
	"stockify/internal/services"
	"stockify/internal/store"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)
	stockStore := store.NewStockStore(db)
	stockService := services.NewStockService(stockStore)
	recommendationService := services.NewRecommendationService(stockStore)
	router := api.NewRouter(stockService, recommendationService)

	log.Printf("Starting server on port %s\n", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
