package models

type Broker struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Cash float64 `json:"cash"`
}

type BrokerDetails struct {
	Broker
	Stocks []StockHolding `json:"stocks"`
}

type StockHolding struct {
	Id         int     `json:"stock_id"`
	Ticker     string  `json:"ticker"`
	Lot        int     `json:"lot"`
	AvgPrice   float64 `json:"avg_price"`
	BrokerName string  `json:"broker_name"`
	BrokerID   int     `json:"broker_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type Portfolio struct {
	ID       int             `json:"id"`
	Brokers  []BrokerDetails `json:"brokers"`
	Holdings []StockHolding  `json:"holdings"`
}

type StockPrice struct {
	Ticker        string  `json:"ticker"`
	Price         float64 `json:"price"`
	PreviousClose float64 `json:"previous_close"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Currency      string  `json:"currency"`
}

func (StockHolding) TableName() string {
	return "stocks"
}

func (Broker) TableName() string {
	return "brokers"
}
