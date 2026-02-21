package contract

import (
	"github.com/marchelrn/stock_api/models"
)

type Repository struct{
	Stocks StocksRepository
}

type StocksRepository interface {
	GetStocks(ticker []string) (*models.StockHolding, error)
}

type BrokerRepository interface {
	GetBrokerDetails(brokerID int) (*models.BrokerDetails, error)
}