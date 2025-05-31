package api

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"stockify/internal/services"
	"stockify/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type StockHandler struct {
	stockService          *services.StockService
	recommendationService *services.RecommendationService
}

func NewStockHandler(ss *services.StockService, rs *services.RecommendationService) *StockHandler {
	return &StockHandler{stockService: ss, recommendationService: rs}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error interno del servidor al generar JSON"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *StockHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	params := store.GetStocksParams{
		Search:    queryParams.Get("search"),
		SortBy:    queryParams.Get("sortBy"),
		SortOrder: queryParams.Get("sortOrder"),
		Page:      page,
		PageSize:  pageSize,
	}

	stocks, totalItems, err := h.stockService.ListStocks(params)
	if err != nil {
		log.Printf("Error en ListStocks service: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Falló la obtención de acciones")
		return
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	response := map[string]interface{}{
		"stocks":     stocks,
		"totalItems": totalItems,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": totalPages,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (h *StockHandler) GetStockByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "Parámetro ticker es requerido")
		return
	}

	stock, err := h.stockService.GetStockByTicker(ticker)
	if err != nil {
		log.Printf("Error en GetStockByTicker service para %s: %v", ticker, err)
		respondWithError(w, http.StatusInternalServerError, "Falló la obtención del stock")
		return
	}
	if stock == nil {
		respondWithError(w, http.StatusNotFound, "Stock no encontrado")
		return
	}
	respondWithJSON(w, http.StatusOK, stock)
}

func (h *StockHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := h.recommendationService.GetRecommendations()
	if err != nil {
		log.Printf("Error en GetRecommendations service: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Falló la obtención de recomendaciones")
		return
	}
	response := map[string]interface{}{
		"recommendations": recommendations,
	}
	respondWithJSON(w, http.StatusOK, response)
}
