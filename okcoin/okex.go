package okcoin

import (
	"net/http"
	"encoding/json"
	"strconv"
	"fmt"
	"errors"
	"net/url"
	"strings"
	. "gexch/common"
)

const (
	OKExName         = "OKEx"
	OKExBaseUri      = "https://www.okex.com/api/v1/"
	FutTickerUri     = "future_ticker.do?symbol=%s&contract_type=%s"
	FutDepthUri      = "future_depth.do?symbol=%s&contract_type=%s"
	FutTradesUri     = "future_trades.do?symbol=%s&contract_type=%s"
	FutIndexUri      = "future_index.do?symbol=%s"
	FutEstimatedUri  = "future_estimated_price.do?symbol=%s"
	FutHoldAmountUri = "future_hold_amount.do?symbol=%s&contract_type=%s"
	FutUserInfoUri   = "future_userinfo.do"
	FutPositionUri   = "future_position.do"
	FutOrderInfo     = "future_order_info.do"
	FutOrdersInfo    = "future_orders_info.do"
	FutTrade         = "future_trade.do"
	FutCancel        = "future_cancel.do"
)

const (
	Weekly    = "this_week"
	BiWeekly  = "next_week"
	Quarterly = "quarter"
)

type OKEx struct {
	client    *http.Client
	apiKey    string
	secretKey string
	Name      string
	BaseUri   string
}

type FutureInfo struct {
	ContractName string
	OpenInterest float64
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

func (ctx *OKEx) GetFutDepth(currency CurrencyPair, contract string, size int) (*Depth, error) {
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
		for i, omap := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = omap.(float64)
			case 1:
				dr.Amount = omap.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}

	for _, v := range dat["bids"].([]interface{}) {
		var dr DepthRecord
		for i, omap := range v.([]interface{}) {
			switch i {
			case 0:
				dr.Price = omap.(float64)
			case 1:
				dr.Amount = omap.(float64)
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
	return dat["future_index"].(float64), nil
}

func (ctx *OKEx) GetFutEstimatedPrice(currency CurrencyPair) (float64, error) {
	url := fmt.Sprintf(ctx.BaseUri+FutEstimatedUri, ExPairSymbol[currency])
	dat, err := HttpGet(ctx.client, url)
	if err != nil {
		return 0, err
	}
	return dat["forecast_price"].(float64), nil
}

func (ctx *OKEx) GetFutureInfo(currency CurrencyPair, contract string) (*FutureInfo, error) {
	url := fmt.Sprintf(ctx.BaseUri+FutHoldAmountUri, ExPairSymbol[currency], contract)
	dat, err := HttpGet2(ctx.client, url)
	if err != nil {
		return nil, err
	}

	rsp := struct {
		Amount       float64
		ContractName string `json:"contract_name"`
	}{}
	err = json.Unmarshal(dat, &rsp)
	if err != nil {
		return nil, err
	}
	return &FutureInfo{ContractName: rsp.ContractName, OpenInterest: rsp.Amount}, nil
}

func (ctx *OKEx) GetFutAccount() (*FutureAccount, error) {
	postData := url.Values{}
	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
	if err != nil {
		return nil, err
	}

	dat, err := HttpPostForm(ctx.client, ctx.BaseUri+FutUserInfoUri, postData)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(dat))

	rsp := struct {
		Info struct {
			Btc map[string]float64 `json:btc`
			Ltc map[string]float64 `json:ltc`
		} `json:info`
		Result bool `json:"result,bool"`
	}{}

	err = json.Unmarshal(dat, &rsp)
	if err != nil {
		return nil, err
	}
	if !rsp.Result {
		return nil, errors.New(string(dat))
	}

	account := new(FutureAccount)
	account.FutureSubAccounts = make(map[Currency]FutureSubAccount, 2)

	btc := rsp.Info.Btc
	ltc := rsp.Info.Ltc

	account.FutureSubAccounts[BTC] = FutureSubAccount{BTC, btc["account_rights"], btc["keep_deposit"], btc["profit_real"], btc["profit_unreal"], btc["risk_rate"]}
	account.FutureSubAccounts[LTC] = FutureSubAccount{LTC, ltc["account_rights"], ltc["keep_deposit"], ltc["profit_real"], ltc["profit_unreal"], ltc["risk_rate"]}

	return account, nil
}

func (ctx *OKEx) GetFutPosition(currency CurrencyPair, contract string) ([]FuturePosition, error) {
	postData := url.Values{}
	postData.Set("symbol", ExPairSymbol[currency])
	postData.Set("contract_type", contract)

	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
	if err != nil {
		return nil, err
	}

	dat, err := HttpPostForm(ctx.client, ctx.BaseUri+FutPositionUri, postData)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dat))

	rsp := struct {
		ForceLiquPirce float64 `json:"force_liqu_pirce,float64"`
		Result         bool `json:"result,bool"`
	}{}

	err = json.Unmarshal(dat, &rsp)
	if err != nil {
		return nil, err
	}
	if !rsp.Result {
		return nil, errors.New(string(dat))
	}

	var pos []FuturePosition
	return pos, nil
}

func (ctx *OKEx) GetFutOrders(currency CurrencyPair, contract string, orderIds []string) ([]FutureOrder, error) {
	postData := url.Values{}
	postData.Set("order_id", strings.Join(orderIds, ","))
	postData.Set("symbol", ExPairSymbol[currency])
	postData.Set("contract_type", contract)

	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
	if err != nil {
		return nil, err
	}

	dat, err := HttpPostForm(ctx.client, ctx.BaseUri+FutOrdersInfo, postData)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dat))

	var respMap map[string]interface{}
	err = json.Unmarshal(dat, &respMap)
	if err != nil {
		return nil, err
	}

	if !respMap["result"].(bool) {
		return nil, errors.New(string(dat))
	}

	var orders []FutureOrder
	for _, v := range respMap["orders"].([]interface{}) {
		order := fillFutureOrder(currency, v.(map[string]interface{}))
		orders = append(orders, order)
	}
	return orders, nil

}

func fillFutureOrder(currency CurrencyPair, omap map[string]interface{}) FutureOrder {
	var order FutureOrder
	order.OrderID = strconv.FormatFloat(omap["order_id"].(float64), 'f', -1, 64)
	order.Amount = omap["amount"].(float64)
	order.Price = omap["price"].(float64)
	order.AvgPrice = omap["price_avg"].(float64)
	order.DealAmount = omap["deal_amount"].(float64)
	order.Fee = omap["fee"].(float64)
	order.OType = int(omap["type"].(float64))
	order.OrderTime = int64(omap["create_date"].(float64))
	order.LeverRate = int(omap["lever_rate"].(float64))
	order.ContractName = omap["contract_name"].(string)
	order.Currency = currency

	switch s := int(omap["status"].(float64)); s {
	case 0:
		order.Status = ORDER_UNFINISH
	case 1:
		order.Status = ORDER_PART_FINISH
	case 2:
		order.Status = ORDER_FINISH
	case 4:
		order.Status = ORDER_CANCEL_ING
	case -1:
		order.Status = ORDER_CANCEL

	}
	return order
}

func (ctx *OKEx) SendFutOrder(ccy CurrencyPair, contract string, price, amount float64, openType, matchPrice, leverRate int) (*FutureOrder, error) {
	postData := url.Values{}
	postData.Set("symbol", ExPairSymbol[ccy])
	postData.Set("price", fmt.Sprint(price))
	postData.Set("contract_type", contract)
	postData.Set("amount", fmt.Sprint(amount))
	postData.Set("type", strconv.Itoa(openType))
	postData.Set("lever_rate", strconv.Itoa(leverRate))
	postData.Set("match_price", strconv.Itoa(matchPrice))

	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
	if err != nil {
		return nil, err
	}

	dat, err := HttpPostForm(ctx.client, ctx.BaseUri+FutTrade, postData)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(dat))

	rsp := struct {
		OrderID int64 `json:"order_id,int64"`
		Result  bool `json:"result,bool"`
	}{}
	err = json.Unmarshal(dat, &rsp)
	if err != nil {
		return nil, err
	}

	if !rsp.Result {
		return nil, errors.New(string(dat))
	}

	fut := new(FutureOrder)
	fut.OrderID = strconv.FormatInt(rsp.OrderID, 10)
	fut.Price = price
	fut.Amount = amount
	fut.Currency = ccy
	fut.OType = openType
	fut.LeverRate = leverRate
	fut.Status = ORDER_UNFINISH
	return fut, nil
}

func (ctx *OKEx) CancelFutOrder(ccy CurrencyPair, contract, orderId string) (bool, error) {
	postData := url.Values{}
	postData.Set("symbol", ExPairSymbol[ccy])
	postData.Set("order_id", orderId)
	postData.Set("contract_type", contract)

	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
	if err != nil {
		return false, err
	}

	dat, err := HttpPostForm(ctx.client, ctx.BaseUri+FutCancel, postData)
	if err != nil {
		return false, err
	}
	fmt.Println(string(dat))

	respMap := make(map[string]interface{})
	err = json.Unmarshal(dat, &respMap)
	if err != nil {
		return false, err
	}

	if respMap["result"] != nil && !respMap["result"].(bool) {
		return false, errors.New(string(dat))
	}

	return true, nil
}

func (ctx *OKEx) GetFutOpenOrders(ccy CurrencyPair, contract string) ([]FutureOrder, error) {
	postData := url.Values{}
	postData.Set("order_id", "-1")
	postData.Set("contract_type", contract)
	postData.Set("symbol", ExPairSymbol[ccy])
	postData.Set("status", "1")
	postData.Set("current_page", "1")
	postData.Set("page_length", "100")

	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
	if err != nil {
		return nil, err
	}

	dat, err := HttpPostForm(ctx.client, ctx.BaseUri+FutOrderInfo, postData)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dat))

	var respMap map[string]interface{}
	err = json.Unmarshal(dat, &respMap)
	if err != nil {
		return nil, err
	}

	if !respMap["result"].(bool) {
		return nil, errors.New(string(dat))
	}

	var orders []FutureOrder
	for _, v := range respMap["orders"].([]interface{}) {
		order := fillFutureOrder(ccy, v.(map[string]interface{}))
		orders = append(orders, order)
	}
	return orders, nil
}
