package huobi

type HuobiTicker struct {
	High float64
	Low  float64
	Last float64
	Vol  float64
	Buy  float64
	Sell float64
}

type HuobiTickerResponse struct {
	Time   string
	Ticker HuobiTicker
}

type HuobiOrderbook struct {
	ID     float64
	TS     float64
	Bids   [][]float64 `json:"bids"`
	Asks   [][]float64 `json:"asks"`
	Symbol string      `json:"string"`
}

type _AccountInfo struct {
	Total        float64 `json:"total,string"`
	NetAsset     float64 `json:"net_asset,string"`
	AvailableBtc float64 `json:"available_btc_display,string"`
	LoanBtc      float64 `json:"loan_btc_display,string"`
	FrozenBtc    float64 `json:"frozen_btc_display,string"`
	AvailableLtc float64 `json:"available_ltc_display,string"`
	LoanLtc      float64 `json:"loan_ltc_display,string"`
	FrozenLtc    float64 `json:"frozen_ltc_display,string"`
	AvailableCny float64 `json:"available_cny_display,string"`
	LoanCny      float64 `json:"loan_cny_display,string"`
	FrozenCny    float64 `json:"frozen_cny_display,string"`
}
