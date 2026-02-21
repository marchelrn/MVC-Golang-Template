package service

import (
	"github.com/marchelrn/stock_api/contract"
	"github.com/marchelrn/stock_api/dto"
)

type HealthService struct {
	health contract.HealthRepository
}

func ImplHealthService(repo contract.HealthRepository) contract.HealthService {
	return &HealthService{health: repo}
}

func (s *HealthService) GetStatus() *dto.HealthResponse {
	return &dto.HealthResponse{
		Status: "healthy",
		Message: "Service is healthy",
	}
}