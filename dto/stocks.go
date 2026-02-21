package dto

type StocksRequest struct {
	Ticker string `json:"ticker" binding:"required"`
	Lot int `json:"lot" binding:"required"`
	AvgPrice float64 `json:"avg_price" binding:"required"`
	BrokerID int `json:"broker_id" binding:"required"`
	BrokerName string `json:"broker_name" binding:"required"`
}

type StocksResponse struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message"`
	StocksData StocksData `json:"stocks_data"`
}

type StocksData struct {
	Ticker string `json:"ticker"`
	Lot int `json:"lot"`
	AvgPrice float64 `json:"avg_price"`
	BrokerID int `json:"broker_id"`
	BrokerName string `json:"broker_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}