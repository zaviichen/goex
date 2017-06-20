package okcoin

import (
	"github.com/zaviichen/gexch/common"
	"net/url"
	"strconv"
	"errors"
)

const (
	OKCnName                 = "OKCn"
	OKComName                = "OKCom"
	OKCnBaseUri              = "https://www.okcoin.cn/api/v1/"
	OKComBaseUri             = "https://www.okcoin.com/api/v1/"
	OKCOIN_TICKER            = "ticker.do"
	OKCOIN_DEPTH             = "depth.do"
	OKCOIN_TRADES            = "trades.do"
	OKCOIN_KLINE             = "kline.do"
	OKCOIN_USERINFO          = "userinfo.do"
	OKCOIN_TRADE             = "trade.do"
	OKCOIN_TRADE_HISTORY     = "trade_history.do"
	OKCOIN_TRADE_BATCH       = "batch_trade.do"
	OKCOIN_ORDER_CANCEL      = "cancel_order.do"
	OKCOIN_ORDER_INFO        = "order_info.do"
	OKCOIN_ORDERS_INFO       = "orders_info.do"
	OKCOIN_ORDER_HISTORY     = "order_history.do"
	OKCOIN_WITHDRAW          = "withdraw.do"
	OKCOIN_WITHDRAW_CANCEL   = "cancel_withdraw.do"
	OKCOIN_WITHDRAW_INFO     = "withdraw_info.do"
	OKCOIN_ORDER_FEE         = "order_fee.do"
	OKCOIN_LEND_DEPTH        = "lend_depth.do"
	OKCOIN_BORROWS_INFO      = "borrows_info.do"
	OKCOIN_BORROW_MONEY      = "borrow_money.do"
	OKCOIN_BORROW_CANCEL     = "cancel_borrow.do"
	OKCOIN_BORROW_ORDER_INFO = "borrow_order_info.do"
	OKCOIN_REPAYMENT         = "repayment.do"
	OKCOIN_UNREPAYMENTS_INFO = "unrepayments_info.do"
	OKCOIN_ACCOUNT_RECORDS   = "account_records.do"
)

type OKSpot struct {
	common.ExchangeBase
}

func NewOKCn(api string, secret string) *OKSpot {
	ex := new(OKSpot)
	ex.Name = OKCnName
	ex.APIUrl = OKCnBaseUri
	ex.Enabled = true
	ex.APIKey = api
	ex.APISecret = secret
	ex.TakerFee = 0.0020
	ex.MakerFee = 0.0020
	ex.Verbose = true
	return ex
}

func NewOKCom(api string, secret string) *OKSpot {
	ex := new(OKSpot)
	ex.Name = OKComName
	ex.APIUrl = OKComBaseUri
	ex.Enabled = true
	ex.APIKey = api
	ex.APISecret = secret
	ex.TakerFee = 0.0020
	ex.MakerFee = 0.0020
	ex.Verbose = true
	return ex
}

func (o *OKSpot) GetTicker(symbol string) (OKCoinTicker, error) {
	resp := OKCoinTickerResponse{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	path := common.EncodeURLValues(o.APIUrl+OKCOIN_TICKER, vals)
	err := common.SendHTTPGetRequest(path, true, &resp)
	if err != nil {
		return OKCoinTicker{}, err
	}
	return resp.Ticker, nil
}

func (o *OKSpot) GetOrderBook(symbol string, size int64, merge bool) (OKCoinOrderbook, error) {
	resp := OKCoinOrderbook{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	if size != 0 {
		vals.Set("size", strconv.FormatInt(size, 10))
	}
	if merge {
		vals.Set("merge", "1")
	}

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_DEPTH, vals)
	err := common.SendHTTPGetRequest(path, true, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (o *OKSpot) GetTrades(symbol string, since int64) ([]OKCoinTrades, error) {
	result := []OKCoinTrades{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	if since != 0 {
		vals.Set("since", strconv.FormatInt(since, 10))
	}

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_TRADES, vals)
	err := common.SendHTTPGetRequest(path, true, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o *OKSpot) GetKline(symbol, klineType string, size, since int64) ([]interface{}, error) {
	resp := []interface{}{}
	vals := url.Values{}
	vals.Set("symbol", symbol)
	vals.Set("type", klineType)

	if size != 0 {
		vals.Set("size", strconv.FormatInt(size, 10))
	}

	if since != 0 {
		vals.Set("since", strconv.FormatInt(since, 10))
	}

	path := common.EncodeURLValues(o.APIUrl+OKCOIN_KLINE, vals)
	err := common.SendHTTPGetRequest(path, true, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OKSpot) GetUserInfo() (OKCoinUserInfo, error) {
	result := OKCoinUserInfo{}
	err := o.PostRequest(OKCOIN_USERINFO, url.Values{}, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (o *OKSpot) Trade(symbol, orderType string, price, amount float64, ) (int64, error) {
	v := url.Values{}
	v.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	v.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	v.Set("symbol", symbol)
	v.Set("type", orderType)

	result := struct {
		Result  bool  `json:"result"`
		OrderID int64 `json:"order_id"`
	}{}
	err := o.PostRequest(OKCOIN_TRADE, v, &result)

	if err != nil {
		return 0, err
	}

	if !result.Result {
		return 0, errors.New("Unable to place order.")
	}
	return result.OrderID, nil
}

func (o *OKSpot) GetTradeHistory(symbol string, since int64) ([]OKCoinTrades, error) {
	result := []OKCoinTrades{}
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("since", strconv.FormatInt(since, 10))

	err := o.PostRequest(OKCOIN_TRADE_HISTORY, v, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o *OKSpot) CancelOrder(symbol string, orderID int64) (int64, error) {
	v := url.Values{}
	v.Set("order_id", strconv.FormatInt(orderID, 10))
	v.Set("symbol", symbol)

	resp := struct {
		OrderID int64 `json:"order_id"`
		Result  bool `json:"result"`
	}{}
	err := o.PostRequest(OKCOIN_ORDER_CANCEL, v, &resp)

	if err != nil {
		return resp.OrderID, err
	}
	return resp.OrderID, nil
}

func (o *OKSpot) GetOrderInfo(symbol string, orderID int64) ([]OKCoinOrderInfo, error) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("order_id", strconv.FormatInt(orderID, 10))

	result := struct {
		Result bool              `json:"result"`
		Orders []OKCoinOrderInfo `json:"orders"`
	}{}
	err := o.PostRequest(OKCOIN_ORDER_INFO, v, &result)

	if err != nil {
		return nil, err
	}

	if result.Result != true {
		return nil, errors.New("Unable to retrieve order info.")
	}
	return result.Orders, nil
}

func (o *OKSpot) GetOrderHistory(pageLength, currentPage int64, status, symbol string) (OKCoinOrderHistory, error) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("status", status)
	v.Set("current_page", strconv.FormatInt(currentPage, 10))
	v.Set("page_length", strconv.FormatInt(pageLength, 10))

	result := OKCoinOrderHistory{}
	err := o.PostRequest(OKCOIN_ORDER_HISTORY, v, &result)

	if err != nil {
		return result, err
	}
	return result, nil
}
