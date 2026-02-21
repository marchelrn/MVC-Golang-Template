package service

import (
	"net/http"

	"github.com/marchelrn/stock_api/contract"
	"github.com/marchelrn/stock_api/dto"
)

type BrokersService struct {
	BrokersRepository contract.BrokerRepository
}

func ImplBrokersService(repo contract.BrokerRepository) contract.BrokersService {
	return &BrokersService{BrokersRepository: repo}
}

func (s *BrokersService) GetBrokersDetails(brokerName []string) (*dto.BrokersResponse, error) {
	brokers, err := s.BrokersRepository.GetBrokerDetails(brokerName)
	if err != nil {
		return nil, err
	}

	var brokerList []dto.Broker

	for _, broker := range brokers {
		var stocksData []dto.StocksData
		for _, stock := range broker.Stocks {
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

		brokerList = append(brokerList, dto.Broker{
			Id:         broker.ID,
			Name:       broker.Name,
			Cash:       broker.Cash,
			StocksData: stocksData,
		})
	}

	response := &dto.BrokersResponse{
		StatusCode: http.StatusOK,
		Message:    "Success received Brokers Details",
		BrokersData: dto.BrokersData{
			Brokers: brokerList,
		},
	}
	return response, nil
}
