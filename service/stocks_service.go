package service

import (
	"errors"
	"log"
	"net/http"

	"github.com/marchelrn/stock_api/contract"
	"github.com/marchelrn/stock_api/dto"
	errs "github.com/marchelrn/stock_api/pkg/error"
	"gorm.io/gorm"
)

type StocksService struct {
	StocksRepository contract.StocksRepository
}

func ImplStocksService(repo contract.StocksRepository) contract.StocksService {
	return &StocksService{StocksRepository: repo}
}

func (s *StocksService) GetStocks(stockTicker []string) (*dto.StocksResponse, error ){
	stock, err := s.StocksRepository.GetStocks(stockTicker)
	if err != nil {
		log.Printf("error", err)
		if errors.Is(err ,gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("Stock not found")
		}
		return nil, errs.InternalServerError("Failed to retrieve stock data")
	}
	response := &dto.StocksResponse{
		StatusCode: http.StatusOK,
		Message: "Success recieved Stocks",
		StocksData: dto.StocksData{
			Ticker: stock.Ticker,
			Lot: stock.Lot,
			AvgPrice: stock.AvgPrice,
			BrokerID: stock.BrokerID,
			BrokerName: stock.BrokerName,
			CreatedAt: stock.CreatedAt,
			UpdatedAt: stock.UpdatedAt,
		},
	}
	return response, nil
}