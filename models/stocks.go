package models

type Broker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cash float64 `json:"cash"`
}

type BrokerDetails struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cash float64 `json:"cash"`
	Stocks []StockHolding `json:"stocks"`
	Broker
}

type StockHolding struct {
	Id int `json:"stock_id"`
	Ticker string `json:"ticker"`
	Lot int `json:"lot"`
	AvgPrice float64 `json:"avg_price"`
	BrokerName string `json:"broker_name"`
	BrokerID int `json:"broker_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (StockHolding) TableName() string {
	return "stocks"
}

type Portfolio struct {
	ID int `json:"id"`
	Brokers []BrokerDetails
	Holdings []StockHolding
}

type StockPrice struct {
	Ticker        string  `json:"ticker"`
	Price         float64 `json:"price"`
	PreviousClose float64 `json:"previous_close"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Currency      string  `json:"currency"`
}