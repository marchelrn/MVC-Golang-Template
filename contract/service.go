package contract

import "github.com/marchelrn/stock_api/dto"

type Service struct {
	Stocks  StocksService
	Brokers BrokersService
}

type StocksService interface {
	GetStocks(stockTicker []string) (*dto.StocksResponse, error)
}

type BrokersService interface {
	GetBrokersDetails(brokersName []string) (*dto.BrokersResponse, error)
}
