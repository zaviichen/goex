package types

import "time"

type MetaInfo struct {
	Exchange   string
	Source     string
	Symbol     string
	ServerTime time.Time
	LocalTime  time.Time
}

type Ticker struct {
	MetaInfo
	Last   float64
	Buy    float64
	Sell   float64
	High   float64
	Low    float64
	Volume float64
}

type OrderBook struct {
	MetaInfo
	AskPrices []float64
	AskSizes  []float64
	BidPrices []float64
	BidSizes  []float64
}

type Trade struct {
	MetaInfo
}