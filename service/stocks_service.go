package service

import (
	"errors"
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

func (s *StocksService) GetStocks(stockTicker []string) (*dto.StocksResponse, error) {
	stocks, err := s.StocksRepository.GetStocks(stockTicker)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("Stock not found")
		}
		return nil, errs.InternalServerError("Failed to retrieve stock data")
	}

	if len(stocks) == 0 {
		return nil, errs.NotFound("Stock not found")
	}

	stocksData := make([]dto.StocksData, 0, len(stocks))
	for _, stock := range stocks {
		stocksData = append(stocksData, dto.StocksData{
			Ticker:     stock.Ticker,
			Lot:        stock.Lot,
			AvgPrice:   stock.AvgPrice,
			BrokerID:   stock.BrokerID,
			BrokerName: stock.BrokerName,
			CreatedAt:  stock.CreatedAt,
			UpdatedAt:  stock.UpdatedAt,
		})
	}

	response := &dto.StocksResponse{
		StatusCode: http.StatusOK,
		Message:    "Success received stocks details",
		StocksData: stocksData,
	}
	return response, nil
}
