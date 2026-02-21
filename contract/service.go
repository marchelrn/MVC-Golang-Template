package contract

import "github.com/marchelrn/stock_api/dto"



type Service struct {
	Health HealthService
}

type HealthService interface {
	GetStatus() *dto.HealthResponse
}