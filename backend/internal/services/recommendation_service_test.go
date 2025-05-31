package services_test

import (
	"errors"
	"stockify/internal/core"
	"stockify/internal/services"
	"stockify/internal/store"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func float64Ptr(v float64) *float64 {
	return &v
}

func TestRecommendationService_GetRecommendations(t *testing.T) {
	mockStore := new(MockStockStore)
	recommendationService := services.NewRecommendationService(mockStore)
	now := time.Now()
	threeMonthsAgo := now.AddDate(0, -3, 0)
	sixMonthsAgo := now.AddDate(0, -6, 0)

	testStocks := []core.Stock{
		{Ticker: "GOOD1", Company: "Good Co One", RatingTo: "Buy", TargetTo: float64Ptr(100.0), Time: now, Action: "upgraded by", Brokerage: "Broker A"},
		{Ticker: "GOOD2", Company: "Good Co Two", RatingTo: "Strong Buy", TargetTo: float64Ptr(150.0), TargetFrom: float64Ptr(120.0), Time: threeMonthsAgo.Add(time.Hour), Action: "target raised by", Brokerage: "Broker B"},
		{Ticker: "HOLD1", Company: "Hold Inc", RatingTo: "Hold", TargetTo: float64Ptr(50.0), Time: now},
		{Ticker: "OLD_BUY", Company: "Old Buy LLC", RatingTo: "Buy", TargetTo: float64Ptr(80.0), Time: sixMonthsAgo},
	}

	expectedParams := store.GetStocksParams{SortBy: "time", SortOrder: "desc", PageSize: 500, Page: 1}
	mockStore.On("GetStocks", expectedParams).Return(testStocks, int64(len(testStocks)), nil).Once()
	recommendations, err := recommendationService.GetRecommendations()

	assert.NoError(t, err)
	assert.NotNil(t, recommendations)
	assert.Len(t, recommendations, 3)

	if len(recommendations) >= 1 {
		assert.Equal(t, "GOOD1", recommendations[0].Ticker)
		assert.Greater(t, recommendations[0].Score, float64(100))
		assert.NotEmpty(t, recommendations[0].Reasons)

		var foundUpgradeReason bool

		for _, reason := range recommendations[0].Reasons {
			if reason.Type == services.ReasonTypeBrokerUpgrade {
				foundUpgradeReason = true
				assert.Contains(t, reason.Details, "Broker A")
			}
		}
		assert.True(t, foundUpgradeReason, "Debería tener una razón de tipo BrokerUpgrade")
	}

	if len(recommendations) >= 2 {
		assert.Equal(t, "GOOD2", recommendations[1].Ticker)
	}

	if len(recommendations) >= 3 {
		assert.Equal(t, "OLD_BUY", recommendations[2].Ticker)
		var foundOldRecentReason bool
		for _, reason := range recommendations[2].Reasons {
			if reason.Type == services.ReasonTypeRecentEvent {
				foundOldRecentReason = true
			}
		}
		assert.False(t, foundOldRecentReason, "OLD_BUY no debería tener la razón de evento reciente")
	}

	mockStore.AssertExpectations(t)
}

func TestRecommendationService_GetRecommendations_StoreError(t *testing.T) {
	mockStore := new(MockStockStore)
	recommendationService := services.NewRecommendationService(mockStore)
	expectedError := errors.New("store failed")
	mockStore.On("GetStocks", mock.AnythingOfType("store.GetStocksParams")).Return(nil, int64(0), expectedError)
	recommendations, err := recommendationService.GetRecommendations()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, recommendations)
	mockStore.AssertExpectations(t)
}
