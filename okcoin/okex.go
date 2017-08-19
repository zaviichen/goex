package okcoin

import (
	"strconv"
	"net/url"
	"github.com/zaviichen/gexch/common"
	"log"
	"github.com/gorilla/websocket"
)

const (
	OKExName                       = "OKEx"
	OKExBaseUri                    = "https://www.okex.com/api/v1/"
	OKCOIN_FUTURES_TICKER          = "future_ticker.do"
	OKCOIN_FUTURES_DEPTH           = "future_depth.do"
	OKCOIN_FUTURES_TRADES          = "future_trades.do"
	OKCOIN_FUTURES_INDEX           = "future_index.do"
	OKCOIN_EXCHANGE_RATE           = "exchange_rate.do"
	OKCOIN_FUTURES_ESTIMATED_PRICE = "future_estimated_price.do"
	OKCOIN_FUTURES_KLINE           = "future_kline.do"
	OKCOIN_FUTURES_HOLD_AMOUNT     = "future_hold_amount.do"
	OKCOIN_FUTURES_USERINFO        = "future_userinfo.do"
	OKCOIN_FUTURES_POSITION        = "future_position.do"
	OKCOIN_FUTURES_TRADE           = "future_trade.do"
	OKCOIN_FUTURES_TRADE_HISTORY   = "future_trades_history.do"
	OKCOIN_FUTURES_TRADE_BATCH     = "future_batch_trade.do"
	OKCOIN_FUTURES_CANCEL          = "future_cancel.do"
	OKCOIN_FUTURES_ORDER_INFO      = "future_order_info.do"
	OKCOIN_FUTURES_ORDERS_INFO     = "future_orders_info.do"
	OKCOIN_FUTURES_USERINFO_4FIX   = "future_userinfo_4fix.do"
	OKCOIN_FUTURES_POSITION_4FIX   = "future_position_4fix.do"
	OKCOIN_FUTURES_EXPLOSIVE       = "future_explosive.do"
	OKCOIN_FUTURES_DEVOLVE         = "future_devolve.do"
)

const (
	Weekly     = "this_week"
	BiWeekly   = "next_week"
	Quarterly  = "quarter"
	OpenLong   = "1"
	OpenShort  = "2"
	CloseLong  = "3"
	CloseShort = "4"
)

type OKEx struct {
	common.ExchangeBase
	OpenFee, CloseFee, DeliveryFee float64
	WebsocketConn                  *websocket.Conn
}

func NewOKEx(api string, secret string) *OKEx {
	ex := new(OKEx)
	ex.Name = OKExName
	ex.APIUrl = OKExBaseUri
	ex.Enabled = true
	ex.APIKey = api
	ex.APISecret = secret
	ex.OpenFee = 0.0003
	ex.CloseFee = 0
	ex.DeliveryFee = 0
	ex.Verbose = true
	return ex
}

func (o *OKEx) GetFuturesTicker(symbol, contractType string) (OKCoinFuturesTicker, error) {
	resp := OKCoinFuturesTickerResponse{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("contract_type", contractType)
	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_TICKER, vals)
	err := common.SendHTTPGetRequest(path, true, &resp)
	if err != nil {
		return OKCoinFuturesTicker{}, err
	}
	return resp.Ticker, nil
}

func (o *OKEx) GetFuturesDepth(symbol, contractType string, size int64, merge bool) (OKCoinOrderbook, error) {
	result := OKCoinOrderbook{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("contract_type", contractType)

	if size != 0 {
		vals.Set("size", strconv.FormatInt(size, 10))
	}
	if merge {
		vals.Set("merge", "1")
	}

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_DEPTH, vals)
	err := common.SendHTTPGetRequest(path, true, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (o *OKEx) GetFuturesTrades(symbol, contractType string) ([]OKCoinFuturesTrades, error) {
	result := []OKCoinFuturesTrades{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("contract_type", contractType)

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_TRADES, vals)
	err := common.SendHTTPGetRequest(path, true, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o *OKEx) GetFuturesIndex(symbol string) (float64, error) {
	result := struct {
		Index float64 `json:"future_index"`
	}{}
	vals := url.Values{}
	vals.Set("symbol", symbol)

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_INDEX, vals)
	err := common.SendHTTPGetRequest(path, true, &result)
	if err != nil {
		return 0, err
	}
	return result.Index, nil
}

func (o *OKEx) GetFuturesExchangeRate() (float64, error) {
	result := struct {
		Rate float64 `json:"rate"`
	}{}
	err := common.SendHTTPGetRequest(o.APIUrl+OKCOIN_EXCHANGE_RATE, true, &result)
	if err != nil {
		return result.Rate, err
	}
	return result.Rate, nil
}

func (o *OKEx) GetFuturesEstimatedPrice(symbol string) (float64, error) {
	result := struct {
		Price float64 `json:"forecast_price"`
	}{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_ESTIMATED_PRICE, vals)
	err := common.SendHTTPGetRequest(path, true, &result)
	if err != nil {
		return result.Price, err
	}
	return result.Price, nil
}

func (o *OKEx) GetFuturesKline(symbol, klineType, contractType string, size, since int64) ([]interface{}, error) {
	resp := []interface{}{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("type", klineType)
	vals.Set("contract_type", contractType)

	if size != 0 {
		vals.Set("size", strconv.FormatInt(size, 10))
	}
	if since != 0 {
		vals.Set("since", strconv.FormatInt(since, 10))
	}

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_KLINE, vals)
	err := common.SendHTTPGetRequest(path, true, &resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *OKEx) GetFuturesHoldAmount(symbol, contractType string) ([]OKCoinFuturesHoldAmount, error) {
	resp := []OKCoinFuturesHoldAmount{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("contract_type", contractType)

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_FUTURES_HOLD_AMOUNT, vals)
	err := common.SendHTTPGetRequest(path, true, &resp)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *OKEx) GetFuturesUserInfo() (OKCoinFuturesUserInfo, error) {
	resp := OKCoinFuturesUserInfo{}
	err := o.PostRequest(OKCOIN_FUTURES_USERINFO, url.Values{}, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (o *OKEx) GetFuturesPosition(symbol, contractType string) (OKCoinFuturesPosition, error) {
	resp := OKCoinFuturesPosition{}
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	err := o.PostRequest(OKCOIN_FUTURES_POSITION, v, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (o *OKEx) FuturesTrade(symbol, contractType, orderType string, price, amount float64, matchPrice, leverage int64) (int64, error) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	v.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	v.Set("type", orderType)
	v.Set("match_price", strconv.FormatInt(matchPrice, 10))
	v.Set("lever_rate", strconv.FormatInt(leverage, 10))

	resp := struct {
		OrderID int64 `json:"order_id"`
		Result  bool `json:"result"`
	}{}
	err := o.PostRequest(OKCOIN_FUTURES_TRADE, v, &resp)
	if err != nil {
		return resp.OrderID, err
	}
	return resp.OrderID, nil
}

func (o *OKEx) CancelFuturesOrder(symbol, contractType string, orderID int64) (int64, error) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("order_id", strconv.FormatInt(orderID, 10))

	resp := struct {
		OrderID int64 `json:"order_id,string"`
		Result  bool `json:"result"`
	}{}
	err := o.PostRequest(OKCOIN_FUTURES_CANCEL, v, &resp)
	if err != nil {
		return resp.OrderID, err
	}
	return resp.OrderID, nil
}

func (o *OKEx) GetFuturesOrdersInfo(symbol, contractType string, orderID int64) ([]OKCoinFuturesOrderInfo, error) {
	v := url.Values{}
	v.Set("order_id", strconv.FormatInt(orderID, 10))
	v.Set("contract_type", contractType)
	v.Set("symbol", symbol)

	resp := struct {
		Orders []OKCoinFuturesOrderInfo `json:"orders"`
		Result bool `json:"result"`
	}{}
	err := o.PostRequest(OKCOIN_FUTURES_ORDERS_INFO, v, &resp)
	if err != nil {
		return resp.Orders, err
	}
	return resp.Orders, nil
}

func (o *OKEx) GetFuturesOrderInfo(orderID, status, currentPage, pageLength int64, symbol, contractType string) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("status", strconv.FormatInt(status, 10))
	v.Set("order_id", strconv.FormatInt(orderID, 10))
	v.Set("current_page", strconv.FormatInt(currentPage, 10))
	v.Set("page_length", strconv.FormatInt(pageLength, 10))

	err := o.PostRequest(OKCOIN_FUTURES_ORDER_INFO, v, nil)
	if err != nil {
		log.Println(err)
	}
}

func (o *OKEx) GetFuturesUserInfo4Fix() {
	v := url.Values{}
	err := o.PostRequest(OKCOIN_FUTURES_USERINFO_4FIX, v, nil)
	if err != nil {
		log.Println(err)
	}
}

func (o *OKEx) GetFuturesUserPosition4Fix(symbol, contractType string) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("type", strconv.FormatInt(1, 10))

	err := o.PostRequest(OKCOIN_FUTURES_POSITION_4FIX, v, nil)

	if err != nil {
		log.Println(err)
	}
}
