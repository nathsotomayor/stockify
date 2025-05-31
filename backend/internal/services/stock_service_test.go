package services_test

import (
	"errors"
	"stockify/internal/core"
	"stockify/internal/services"
	"stockify/internal/store"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStockStore struct {
	mock.Mock
	store.StockStoreInterface
}

func (m *MockStockStore) GetStocks(params store.GetStocksParams) ([]core.Stock, int64, error) {
	args := m.Called(params)

	var stocks []core.Stock

	if arg0 := args.Get(0); arg0 != nil {
		stocks = arg0.([]core.Stock)
	}

	return stocks, args.Get(1).(int64), args.Error(2)
}

func (m *MockStockStore) GetStockByTicker(ticker string) (*core.Stock, error) {
	args := m.Called(ticker)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*core.Stock), args.Error(1)
}

func (m *MockStockStore) CountStocks() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStockStore) GetRawStocksForRecommendation(limit int) ([]core.Stock, error) {
	args := m.Called(limit)

	var stocks []core.Stock

	if arg0 := args.Get(0); arg0 != nil {
		stocks = arg0.([]core.Stock)
	}

	return stocks, args.Error(1)
}

func TestStockService_ListStocks(t *testing.T) {
	mockStore := new(MockStockStore)
	stockService := services.NewStockService(mockStore)
	expectedStocks := []core.Stock{{Ticker: "AAPL"}}
	expectedTotal := int64(1)
	params := store.GetStocksParams{Page: 1, PageSize: 10}

	mockStore.On("GetStocks", params).Return(expectedStocks, expectedTotal, nil)

	stocks, total, err := stockService.ListStocks(params)

	assert.NoError(t, err)
	assert.Equal(t, expectedStocks, stocks)
	assert.Equal(t, expectedTotal, total)
	mockStore.AssertExpectations(t)
}

func TestStockService_GetStockByTicker_Found(t *testing.T) {
	mockStore := new(MockStockStore)
	stockService := services.NewStockService(mockStore)
	expectedStock := &core.Stock{Ticker: "AAPL", Company: "Apple Inc."}
	ticker := "AAPL"

	mockStore.On("GetStockByTicker", ticker).Return(expectedStock, nil)

	stock, err := stockService.GetStockByTicker(ticker)

	assert.NoError(t, err)
	assert.Equal(t, expectedStock, stock)
	mockStore.AssertExpectations(t)
}

func TestStockService_GetStockByTicker_NotFound(t *testing.T) {
	mockStore := new(MockStockStore)
	stockService := services.NewStockService(mockStore)
	ticker := "UNKNOWN"

	mockStore.On("GetStockByTicker", ticker).Return(nil, nil)

	stock, err := stockService.GetStockByTicker(ticker)

	assert.NoError(t, err)
	assert.Nil(t, stock)
	mockStore.AssertExpectations(t)
}

func TestStockService_GetStockByTicker_StoreError(t *testing.T) {
	mockStore := new(MockStockStore)
	stockService := services.NewStockService(mockStore)
	ticker := "ERROR"
	expectedError := errors.New("database error")

	mockStore.On("GetStockByTicker", ticker).Return(nil, expectedError)

	stock, err := stockService.GetStockByTicker(ticker)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, stock)
	mockStore.AssertExpectations(t)
}
