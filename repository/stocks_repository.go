package repository

import (
	"github.com/marchelrn/stock_api/contract"
	"github.com/marchelrn/stock_api/models"

	"gorm.io/gorm"
)

type StocksRepository struct{
	db *gorm.DB
}


func ImplStocksRepository(db *gorm.DB) contract.StocksRepository {
	return &StocksRepository{db: db}
}

func (r *StocksRepository) GetStocks(ticker []string) (*models.StockHolding, error) {
	var stock models.StockHolding
	if err := r.db.Where("ticker IN ?", ticker).First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

