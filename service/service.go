package service

import (
	"github.com/marchelrn/stock_api/contract"
)

func New(repo *contract.Repository) (*contract.Service, error) {
	return &contract.Service{
		Health: ImplHealthService(repo.Health),
	}, nil
}