package repository

import (
	"github.com/marchelrn/stock_api/models"
	"gorm.io/gorm"
)

type BrokersRepository struct {
	db *gorm.DB
}

func ImplBrokersRepository(db *gorm.DB) *BrokersRepository {
	return &BrokersRepository{db: db}
}

func (r *BrokersRepository) GetBrokerDetails(brokerName []string) ([]models.BrokerDetails, error) {
	var brokerDetails []models.BrokerDetails
	if err := r.db.Where("name IN ?", brokerName).Find(&brokerDetails).Error; err != nil {
		return nil, err
	}
	return brokerDetails, nil
}
