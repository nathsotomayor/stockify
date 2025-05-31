package store

import "stockify/internal/core"

type StockStoreInterface interface {
	GetStocks(params GetStocksParams) ([]core.Stock, int64, error)
	GetStockByTicker(ticker string) (*core.Stock, error)
	CountStocks() (int64, error)
	GetRawStocksForRecommendation(limit int) ([]core.Stock, error)
}
