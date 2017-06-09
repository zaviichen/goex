package okcoin

import (
	"net/http"
	"encoding/json"
	"strconv"
	"fmt"
	"errors"
)

const (
	OKExName     = "OKEx"
	OKExBaseUri  = "https://www.okcoin.com/api/v1/"
	FutTickerUri = "future_ticker.do?symbol=%s&contract_type=%s"
	FutDepthUri  = "future_depth.do?symbol=%s&contract_type=%s"
	FutTradesUri = "future_trades.do?symbol=%s&contract_type=%s"
	FutIndexUri = "future_index.do?symbol=%s"
)

const (
	Weekly = "this_week"
	BiWeekly = "next_week"
	Quarterly = "quarter"
)

type OKEx struct {
	client    *http.Client
	apiKey    string
	secretKey string
	Name      string
	BaseUri   string
}

func NewOKEx(client *http.Client, api string, secret string) *OKEx {
	ex := new(OKEx)
	ex.client = client
	ex.apiKey = api
	ex.secretKey = secret
	ex.Name = OKExName
	ex.BaseUri = OKExBaseUri
	return ex
}

func (ctx *OKEx) GetFutTicker(currency CurrencyPair, contract string) (*Ticker, error) {
	url := fmt.Sprintf(ctx.BaseUri+FutTickerUri, ExPairSymbol[currency], contract)
	dat, err := HttpGet2(ctx.client, url)
	if err != nil {
		return nil, err
	}

	rsp := struct {
		Date string
		Ticker struct {
			Last       float64
			High       float64
			Low        float64
			Buy        float64
			Sell       float64
			Vol        float64
			ContractID int `json:"contract_id"`
			UnitAmount int `json:"unit_amount"`
		}
	}{}
	err = json.Unmarshal(dat, &rsp)
	if err != nil {
		return nil, err
	}

	t := new(Ticker)
	t.Date, _ = strconv.ParseUint(rsp.Date, 10, 64)
	t.Buy = rsp.Ticker.Buy
	t.Sell = rsp.Ticker.Sell
	t.High = rsp.Ticker.High
	t.Last = rsp.Ticker.Last
	t.Low = rsp.Ticker.Low
	t.Vol = rsp.Ticker.Vol
	return t, nil
}


func (ctx *OKEx) GetFutDepth(currency CurrencyPair, contract string , size int) (*Depth, error) {
	url := fmt.Sprintf(ctx.BaseUri+FutDepthUri, ExPairSymbol[currency], contract)
	dat, err := HttpGet(ctx.client, url)
	if err != nil {
		return nil, err
	}

	if dat["result"] != nil && !dat["result"].(bool) {
		return nil, errors.New(fmt.Sprintf("%.0f", dat["error_code"].(float64)))
	}

	var depth Depth
	for _, v := range dat["asks"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range dat["bids"].([]interface{}) {
		var dr DepthRecord
		for i, vv := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = vv.(float64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.BidList = append(depth.BidList, dr)
	}

	return &depth, nil
}


func (ctx *OKEx) GetFutTrades(currency CurrencyPair, contract string) ([]Trade, error) {
	url := fmt.Sprintf(ctx.BaseUri+FutTradesUri, ExPairSymbol[currency], contract)
	dat, err := HttpGet2(ctx.client, url)
	if err != nil {
		return nil, err
	}

	type rspTrade struct {
		Tid    int64
		Type   string
		Amount float64
		Price  float64
		Date   int64
		DateMs int64 `json:"date_ms"`
	}
	var rsp []rspTrade
	err = json.Unmarshal(dat, &rsp)
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, v := range rsp {
		trade := new(Trade)
		trade.Tid = v.Tid
		trade.Type = v.Type
		trade.Amount = v.Amount
		trade.Price = v.Price
		trade.Date = v.DateMs
		trades = append(trades, *trade)
	}
	return trades, nil
}

func (ctx *OKEx) GetFutIndex(currency CurrencyPair) (float64, error) {
	url := fmt.Sprintf(ctx.BaseUri+FutIndexUri, ExPairSymbol[currency])
	dat, err := HttpGet(ctx.client, url)
	if err != nil {
		return 0, err
	}

	index := dat["future_index"].(float64)
	return index, nil
}