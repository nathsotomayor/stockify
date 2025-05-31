package services

import (
	"fmt"
	"log"
	"sort"
	"stockify/internal/core"
	"stockify/internal/store"
	"strings"
	"time"
)

type RecommendationService struct {
	stockStore store.StockStoreInterface
}

func NewRecommendationService(ss store.StockStoreInterface) *RecommendationService {
	return &RecommendationService{stockStore: ss}
}

type RecommendationReasonType string

const (
	ReasonTypePositiveRating   RecommendationReasonType = "POSITIVE_RATING"
	ReasonTypeTargetIncreased  RecommendationReasonType = "TARGET_INCREASED"
	ReasonTypeTargetAttractive RecommendationReasonType = "TARGET_ATTRACTIVE"
	ReasonTypeBrokerUpgrade    RecommendationReasonType = "BROKER_UPGRADE"
	ReasonTypeNewCoverage      RecommendationReasonType = "NEW_POSITIVE_COVERAGE"
	ReasonTypeRecentEvent      RecommendationReasonType = "RECENT_EVENT"
)

type RecommendationReason struct {
	Type    RecommendationReasonType `json:"type"`
	Details string                   `json:"details"`
}

type RecommendedStock struct {
	core.Stock
	Reasons []RecommendationReason `json:"reasons"`
	Score   float64                `json:"score"`
}

func (svc *RecommendationService) GetRecommendations() ([]RecommendedStock, error) {
	log.Println("RecommendationService: Iniciando obtención de recomendaciones...")

	allStocks, _, err := svc.stockStore.GetStocks(store.GetStocksParams{
		SortBy:    "time",
		SortOrder: "desc",
		PageSize:  500,
		Page:      1,
	})

	if err != nil {
		log.Printf("RecommendationService: Error obteniendo stocks para recomendaciones: %v", err)
		return nil, err
	}

	log.Printf("RecommendationService: Obtenidos %d stocks para analizar.", len(allStocks))

	var candidates []RecommendedStock

	positiveRatings := map[string]bool{
		"buy":        true,
		"outperform": true,
		"strong buy": true,
		"overweight": true,
		"accumulate": true,
		"add":        true,
		"positive":   true,
	}

	maxAgeForHighConsideration := time.Now().AddDate(0, -3, 0)

	for _, stock := range allStocks {
		var score float64 = 0
		var reasons []RecommendationReason

		if positiveRatings[strings.ToLower(stock.RatingTo)] {
			score += 50
			reasons = append(reasons, RecommendationReason{
				Type:    ReasonTypePositiveRating,
				Details: fmt.Sprintf("Rating positivo: %s", stock.RatingTo),
			})
		} else {
			continue
		}

		if stock.TargetTo != nil {
			score += (*stock.TargetTo / 10)
			if stock.TargetFrom != nil && *stock.TargetTo > *stock.TargetFrom {
				score += 20
				reasons = append(reasons, RecommendationReason{
					Type:    ReasonTypeTargetIncreased,
					Details: fmt.Sprintf("Precio objetivo aumentado de $%.2f a $%.2f", *stock.TargetFrom, *stock.TargetTo),
				})
			} else {
				reasons = append(reasons, RecommendationReason{
					Type:    ReasonTypeTargetAttractive,
					Details: fmt.Sprintf("Precio objetivo atractivo: $%.2f", *stock.TargetTo),
				})
			}
		}

		if strings.Contains(strings.ToLower(stock.Action), "upgraded by") {
			score += 30
			reasons = append(reasons, RecommendationReason{
				Type:    ReasonTypeBrokerUpgrade,
				Details: fmt.Sprintf("Mejorada por %s", stock.Brokerage),
			})
		}

		if strings.Contains(strings.ToLower(stock.Action), "initiated by") && positiveRatings[strings.ToLower(stock.RatingTo)] {
			score += 25
			reasons = append(reasons, RecommendationReason{
				Type:    ReasonTypeNewCoverage,
				Details: fmt.Sprintf("Nueva cobertura (%s) iniciada con rating positivo: %s", stock.Brokerage, stock.RatingTo),
			})
		}

		if stock.Time.After(maxAgeForHighConsideration) {
			score += 15
			reasons = append(reasons, RecommendationReason{
				Type:    ReasonTypeRecentEvent,
				Details: "Evento de rating reciente (últimos 3 meses).",
			})
		}

		if score > 50 {
			candidates = append(candidates, RecommendedStock{
				Stock:   stock,
				Reasons: reasons,
				Score:   score,
			})
		}
	}

	log.Printf("RecommendationService: Encontrados %d candidatos después del scoring.", len(candidates))

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	numRecommendations := 5
	if len(candidates) < numRecommendations {
		numRecommendations = len(candidates)
	}

	finalRecommendations := candidates[:numRecommendations]
	log.Printf("RecommendationService: Devolviendo %d recomendaciones.", len(finalRecommendations))

	return finalRecommendations, nil
}
