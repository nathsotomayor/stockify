package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"stockify/internal/config"
	"stockify/internal/core"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ExtAPIResponse struct {
	Items    []ExtAPIStockItem `json:"items"`
	NextPage *string           `json:"next_page,omitempty"`
}

type ExtAPIStockItem struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

type DataSyncService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewDataSyncService(db *gorm.DB, cfg *config.Config) *DataSyncService {
	return &DataSyncService{db: db, cfg: cfg}
}

func parseMonetaryValue(valueStr string) (*float64, error) {
	if strings.TrimSpace(valueStr) == "" {
		return nil, nil
	}

	cleanedStr := strings.ReplaceAll(valueStr, "$", "")
	cleanedStr = strings.ReplaceAll(cleanedStr, ",", "")
	val, err := strconv.ParseFloat(strings.TrimSpace(cleanedStr), 64)

	if err != nil {
		return nil, fmt.Errorf("no se pudo convertir '%s' a float64: %w", valueStr, err)
	}

	return &val, nil
}

func (s *DataSyncService) RunPopulation() error {
	log.Println("Iniciando tarea de población de la base de datos desde API externa...")

	baseURL := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	apiToken := s.cfg.StockAPIToken
	httpClient := &http.Client{Timeout: 60 * time.Second}
	currentNextPageToken := ""
	pageCount := 1
	totalItemsProcessed := 0

	for {
		var currentApiURL string
		parsedBaseURL, _ := url.Parse(baseURL)
		query := parsedBaseURL.Query()

		if currentNextPageToken != "" {
			query.Set("next_page", currentNextPageToken)
		}

		parsedBaseURL.RawQuery = query.Encode()
		currentApiURL = parsedBaseURL.String()

		log.Printf("Poblando: Obteniendo datos de API externa (Página %d): %s\n", pageCount, currentApiURL)

		req, err := http.NewRequest("GET", currentApiURL, nil)

		if err != nil {
			return fmt.Errorf("poblando: error creando petición HTTP (Página %d): %w", pageCount, err)
		}

		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Accept", "application/json")

		resp, err := httpClient.Do(req)

		if err != nil {
			log.Printf("Poblando: error obteniendo datos de API externa (Página %d): %v. Reintentando en 10 segundos...\n", pageCount, err)
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("Poblando: petición a API externa falló (Página %d) con código %d: %s.\n", pageCount, resp.StatusCode, string(bodyBytes))

			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				log.Println("Poblando: Error de servidor o límite de peticiones. Reintentando en 20 segundos...")
				time.Sleep(20 * time.Second)
				continue
			}

			return fmt.Errorf("poblando: error API no recuperable (código %d)", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return fmt.Errorf("poblando: error leyendo cuerpo de respuesta (Página %d): %w", pageCount, err)
		}

		var apiResponse ExtAPIResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return fmt.Errorf("poblando: error decodificando JSON (Página %d): %w\nCuerpo: %s", pageCount, err, string(body))
		}

		log.Printf("Poblando: Obtenidos %d ítems de API externa (Página %d). Procesando...\n", len(apiResponse.Items), pageCount)
		if len(apiResponse.Items) == 0 && (apiResponse.NextPage == nil || *apiResponse.NextPage == "") {
			log.Println("Poblando: Recibidos 0 ítems y sin token de siguiente página. Asumiendo fin de datos.")
			break
		}

		for _, apiItem := range apiResponse.Items {
			targetFrom, errTFrom := parseMonetaryValue(apiItem.TargetFrom)

			if errTFrom != nil {
				log.Printf("Poblando (Advertencia Ticker %s): TargetFrom ('%s'): %v", apiItem.Ticker, apiItem.TargetFrom, errTFrom)
			}

			targetTo, errTTo := parseMonetaryValue(apiItem.TargetTo)

			if errTTo != nil {
				log.Printf("Poblando (Advertencia Ticker %s): TargetTo ('%s'): %v", apiItem.Ticker, apiItem.TargetTo, errTTo)
			}

			var ratingFromPtr *string

			if trimmedRF := strings.TrimSpace(apiItem.RatingFrom); trimmedRF != "" {
				ratingFromPtr = &trimmedRF
			}

			var parsedTime time.Time

			if strings.TrimSpace(apiItem.Time) != "" {
				var parseErr error
				parsedTime, parseErr = time.Parse(time.RFC3339Nano, apiItem.Time)

				if parseErr != nil {
					log.Printf("Poblando (Advertencia Ticker %s): Time ('%s'): %v. Usando valor zero.", apiItem.Time, apiItem.Ticker, parseErr)
				}
			}

			stockEntry := core.Stock{
				Ticker:     strings.TrimSpace(apiItem.Ticker),
				Company:    strings.TrimSpace(apiItem.Company),
				Brokerage:  strings.TrimSpace(apiItem.Brokerage),
				Action:     strings.TrimSpace(apiItem.Action),
				RatingTo:   strings.TrimSpace(apiItem.RatingTo),
				RatingFrom: ratingFromPtr,
				TargetTo:   targetTo,
				TargetFrom: targetFrom,
				Time:       parsedTime,
			}

			result := s.db.Create(&stockEntry)

			if result.Error != nil {
				log.Printf("Poblando: Error guardando stock para ticker %s en BD: %v\n", stockEntry.Ticker, result.Error)
			} else {
				totalItemsProcessed++
			}
		}

		log.Printf("Poblando: Procesados %d ítems para página %d.\n", len(apiResponse.Items), pageCount)

		if apiResponse.NextPage != nil && strings.TrimSpace(*apiResponse.NextPage) != "" {
			currentNextPageToken = strings.TrimSpace(*apiResponse.NextPage)
			pageCount++
			log.Printf("Poblando: Token de siguiente página encontrado: '%s'. Preparando siguiente página.\n", currentNextPageToken)
		} else {
			log.Println("Poblando: No hay más páginas para obtener.")
			break
		}
	}

	log.Printf("Tarea de población finalizada. Total de ítems procesados (nuevos registros creados o intentados): %d\n", totalItemsProcessed)
	return nil
}
