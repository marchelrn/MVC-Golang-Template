package contract

import "github.com/marchelrn/stock_api/dto"



type Service struct {
	Stocks StocksService
}

type StocksService interface {
	GetStocks(stockTicker []string) (*dto.StocksResponse, error)
}