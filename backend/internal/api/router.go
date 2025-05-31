package api

import (
	"log"
	"net/http"
	"stockify/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(stockService *services.StockService, recommendationService *services.RecommendationService) http.Handler {
	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "Cache-Control", "Pragma", "Expires"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(corsMiddleware.Handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	stockHandler := NewStockHandler(stockService, recommendationService)

	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/stocks", func(stocksRouter chi.Router) {
			stocksRouter.Get("/", stockHandler.GetStocks)
			stocksRouter.Get("/recommendations", stockHandler.GetRecommendations)
			stocksRouter.Get("/{ticker}", stockHandler.GetStockByTicker)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Router configurado con CORS y rutas API.")
	return r
}
