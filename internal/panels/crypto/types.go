package crypto

import "time"

type apiResponse map[string]struct {
	USD       float64 `json:"usd"`
	Change24h float64 `json:"usd_24h_change"`
	MarketCap float64 `json:"usd_market_cap"`
	Volume24h float64 `json:"usd_24h_vol"`
}

type CoinData struct {
	ID        string
	Symbol    string
	Price     float64
	Change24h float64
	MarketCap float64
	Volume24h float64
}

type ResponseMsg struct {
	Coins []CoinData
	Error error
}

type TickMsg time.Time
