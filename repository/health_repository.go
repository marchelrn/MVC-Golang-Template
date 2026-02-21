package repository

import (
	"github.com/marchelrn/stock_api/contract"

	"gorm.io/gorm"
)

type HealthRepository struct{
	db *gorm.DB
}


func ImplHealthRepository(db *gorm.DB) contract.HealthRepository {
	return &HealthRepository{db: db}
}

func (r *HealthRepository) GetStatus() string {
	return "OK"
}

