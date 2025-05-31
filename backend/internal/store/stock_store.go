package store

import (
	"stockify/internal/core"
	"strings"

	"gorm.io/gorm"
)

type StockStore struct {
	db *gorm.DB
}

func NewStockStore(db *gorm.DB) *StockStore {
	return &StockStore{db: db}
}

type GetStocksParams struct {
	Search    string
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

func (s *StockStore) GetStocks(params GetStocksParams) ([]core.Stock, int64, error) {
	var stocks []core.Stock
	var totalItems int64

	query := s.db.Model(&core.Stock{})

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(ticker) LIKE ? OR LOWER(company) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if params.SortBy != "" {
		order := params.SortBy

		if strings.ToLower(params.SortOrder) == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}

		query = query.Order(order)
	} else {
		query = query.Order("time DESC")
	}

	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Limit(params.PageSize).Offset(offset)
	} else if params.PageSize > 0 {
		query = query.Limit(params.PageSize).Offset(0)
	}

	if err := query.Find(&stocks).Error; err != nil {
		return nil, totalItems, err
	}

	return stocks, totalItems, nil
}

func (s *StockStore) GetStockByTicker(ticker string) (*core.Stock, error) {
	var stock core.Stock

	if err := s.db.Where("UPPER(ticker) = ?", strings.ToUpper(ticker)).First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &stock, nil
}

func (s *StockStore) GetRawStocksForRecommendation(limit int) ([]core.Stock, error) {
	var stocks []core.Stock
	query := s.db.Model(&core.Stock{}).Order("time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

func (s *StockStore) CountStocks() (int64, error) {
	var count int64

	if err := s.db.Model(&core.Stock{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
