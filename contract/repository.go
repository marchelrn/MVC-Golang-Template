package contract

import (
	"github.com/marchelrn/stock_api/models"
)

type Repository struct {
	Stocks  StocksRepository
	Brokers BrokerRepository
}

type StocksRepository interface {
	GetStocks(ticker []string) ([]models.StockHolding, error)
}

type BrokerRepository interface {
	GetBrokerDetails(brokerName []string) ([]models.BrokerDetails, error)
}
