package services

import (
	"stockify/internal/core"
	"stockify/internal/store"
)

type StockService struct {
	store store.StockStoreInterface
}

func NewStockService(s store.StockStoreInterface) *StockService {
	return &StockService{store: s}
}

func (svc *StockService) ListStocks(params store.GetStocksParams) ([]core.Stock, int64, error) {
	return svc.store.GetStocks(params)
}

func (svc *StockService) GetStockByTicker(ticker string) (*core.Stock, error) {
	return svc.store.GetStockByTicker(ticker)
}
