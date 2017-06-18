package okcoin

import (
	"strconv"
	"net/url"
	"github.com/zaviichen/gexch/common"
	"log"
	"strings"
	"errors"
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
	OKCOIN_API_URL                 = "https://www.okcoin.com/api/v1/"
	OKCOIN_API_URL_CHINA           = "https://www.okcoin.cn/api/v1/"
	OKCOIN_API_VERSION             = "1"
	OKCOIN_WEBSOCKET_URL           = "wss://real.okcoin.com:10440/websocket/okcoinapi"
	OKCOIN_WEBSOCKET_URL_CHINA     = "wss://real.okcoin.cn:10440/websocket/okcoinapi"
	OKCOIN_TICKER                  = "ticker.do"
	OKCOIN_DEPTH                   = "depth.do"
	OKCOIN_TRADES                  = "trades.do"
	OKCOIN_KLINE                   = "kline.do"
	OKCOIN_USERINFO                = "userinfo.do"
	OKCOIN_TRADE                   = "trade.do"
	OKCOIN_TRADE_HISTORY           = "trade_history.do"
	OKCOIN_TRADE_BATCH             = "batch_trade.do"
	OKCOIN_ORDER_CANCEL            = "cancel_order.do"
	OKCOIN_ORDER_INFO              = "order_info.do"
	OKCOIN_ORDERS_INFO             = "orders_info.do"
	OKCOIN_ORDER_HISTORY           = "order_history.do"
	OKCOIN_WITHDRAW                = "withdraw.do"
	OKCOIN_WITHDRAW_CANCEL         = "cancel_withdraw.do"
	OKCOIN_WITHDRAW_INFO           = "withdraw_info.do"
	OKCOIN_ORDER_FEE               = "order_fee.do"
	OKCOIN_LEND_DEPTH              = "lend_depth.do"
	OKCOIN_BORROWS_INFO            = "borrows_info.do"
	OKCOIN_BORROW_MONEY            = "borrow_money.do"
	OKCOIN_BORROW_CANCEL           = "cancel_borrow.do"
	OKCOIN_BORROW_ORDER_INFO       = "borrow_order_info.do"
	OKCOIN_REPAYMENT               = "repayment.do"
	OKCOIN_UNREPAYMENTS_INFO       = "unrepayments_info.do"
	OKCOIN_ACCOUNT_RECORDS         = "account_records.do"
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

var (
	okcoinDefaultsSet = false
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

type OKCoin struct {
	common.ExchangeBase
	OpenFee, CloseFee, DeliveryFee float64
}

type FutureInfo struct {
	ContractName string
	OpenInterest float64
}

func NewOKEx(api string, secret string) *OKCoin {
	ex := new(OKCoin)
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

func (o *OKCoin) GetFuturesTicker(symbol, contractType string) (OKCoinFuturesTicker, error) {
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

func (o *OKCoin) GetFuturesDepth(symbol, contractType string, size int64, merge bool) (OKCoinOrderbook, error) {
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

func (o *OKCoin) GetFuturesTrades(symbol, contractType string) ([]OKCoinFuturesTrades, error) {
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

func (o *OKCoin) GetFuturesIndex(symbol string) (float64, error) {
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

func (o *OKCoin) GetFuturesExchangeRate() (float64, error) {
	result := struct {
		Rate float64 `json:"rate"`
	}{}
	err := common.SendHTTPGetRequest(o.APIUrl+OKCOIN_EXCHANGE_RATE, true, &result)
	if err != nil {
		return result.Rate, err
	}
	return result.Rate, nil
}

func (o *OKCoin) GetFuturesEstimatedPrice(symbol string) (float64, error) {
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

func (o *OKCoin) GetFuturesKline(symbol, klineType, contractType string, size, since int64) ([]interface{}, error) {
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

func (o *OKCoin) GetFuturesHoldAmount(symbol, contractType string) ([]OKCoinFuturesHoldAmount, error) {
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

func (o *OKCoin) GetFuturesUserInfo() (OKCoinFuturesUserInfo, error) {
	resp := OKCoinFuturesUserInfo{}
	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_USERINFO, url.Values{}, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (o *OKCoin) GetFuturesPosition(symbol, contractType string) (OKCoinFuturesPosition, error) {
	resp := OKCoinFuturesPosition{}
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_POSITION, v, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (o *OKCoin) FuturesTrade(symbol, contractType, orderType string, price, amount float64, matchPrice, leverage int64) (int64, error) {
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
	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_TRADE, v, &resp)
	if err != nil {
		return resp.OrderID, err
	}
	return resp.OrderID, nil
}

func (o *OKCoin) CancelFuturesOrder(symbol, contractType string, orderID int64) (int64, error) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("order_id", strconv.FormatInt(orderID, 10))

	resp := struct {
		OrderID int64 `json:"order_id,string"`
		Result  bool `json:"result"`
	}{}
	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_CANCEL, v, &resp)
	if err != nil {
		return resp.OrderID, err
	}
	return resp.OrderID, nil
}

func (o *OKCoin) GetFuturesOrdersInfo(symbol, contractType string, orderID int64) ([]OKCoinFuturesOrderInfo, error) {
	v := url.Values{}
	v.Set("order_id", strconv.FormatInt(orderID, 10))
	v.Set("contract_type", contractType)
	v.Set("symbol", symbol)

	resp := struct {
		Orders []OKCoinFuturesOrderInfo `json:"orders"`
		Result bool `json:"result"`
	}{}
	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_ORDERS_INFO, v, &resp)
	if err != nil {
		return resp.Orders, err
	}
	return resp.Orders, nil
}

func (o *OKCoin) GetFuturesOrderInfo(orderID, status, currentPage, pageLength int64, symbol, contractType string) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("status", strconv.FormatInt(status, 10))
	v.Set("order_id", strconv.FormatInt(orderID, 10))
	v.Set("current_page", strconv.FormatInt(currentPage, 10))
	v.Set("page_length", strconv.FormatInt(pageLength, 10))

	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_ORDER_INFO, v, nil)
	if err != nil {
		log.Println(err)
	}
}

func (o *OKCoin) GetFuturesUserInfo4Fix() {
	v := url.Values{}
	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_USERINFO_4FIX, v, nil)
	if err != nil {
		log.Println(err)
	}
}

func (o *OKCoin) GetFuturesUserPosition4Fix(symbol, contractType string) {
	v := url.Values{}
	v.Set("symbol", symbol)
	v.Set("contract_type", contractType)
	v.Set("type", strconv.FormatInt(1, 10))

	err := o.SendAuthenticatedHTTPRequest(OKCOIN_FUTURES_POSITION_4FIX, v, nil)

	if err != nil {
		log.Println(err)
	}
}

func (o *OKCoin) SendAuthenticatedHTTPRequest(method string, v url.Values, result interface{}) (err error) {
	v.Set("api_key", o.APIKey)
	hasher := common.GetMD5([]byte(v.Encode() + "&secret_key=" + o.APISecret))
	v.Set("sign", strings.ToUpper(common.HexEncodeToString(hasher)))

	encoded := v.Encode()
	path := o.APIUrl + method

	if o.Verbose {
		log.Printf("Sending POST request to %s with params %s\n", path, encoded)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resp, err := common.SendHTTPRequest("POST", path, headers, strings.NewReader(encoded))

	if err != nil {
		return err
	}

	if o.Verbose {
		log.Printf("Recieved raw: \n%s\n", resp)
	}

	err = common.JSONDecode([]byte(resp), &result)

	if err != nil {
		return errors.New("Unable to JSON Unmarshal response.")
	}

	return nil
}

func RESTErrorMessage(code string) (string) {
	errMsgs := map[string]string{
		"10000": "Required field, can not be null",
		"10001": "Request frequency too high",
		"10002": "System error",
		"10003": "Not in reqest list, please try again later",
		"10004": "IP not allowed to access the resource",
		"10005": "'secretKey' does not exist",
		"10006": "'partner' does not exist",
		"10007": "Signature does not match",
		"10008": "Illegal parameter",
		"10009": "Order does not exist",
		"10010": "Insufficient funds",
		"10011": "Amount too low",
		"10012": "Only btc_usd/btc_cny ltc_usd,ltc_cny supported",
		"10013": "Only support https request",
		"10014": "Order price must be between 0 and 1,000,000",
		"10015": "Order price differs from current market price too much",
		"10016": "Insufficient coins balance",
		"10017": "API authorization error",
		"10018": "Borrow amount less than lower limit [usd/cny:100,btc:0.1,ltc:1]",
		"10019": "Loan agreement not checked",
		"10020": `Rate cannot exceed 1%`,
		"10021": `Rate cannot less than 0.01%`,
		"10023": "Fail to get latest ticker",
		"10024": "Balance not sufficient",
		"10025": "Quota is full, cannot borrow temporarily",
		"10026": "Loan (including reserved loan) and margin cannot be withdrawn",
		"10027": "Cannot withdraw within 24 hrs of authentication information modification",
		"10028": "Withdrawal amount exceeds daily limit",
		"10029": "Account has unpaid loan, please cancel/pay off the loan before withdraw",
		"10031": "Deposits can only be withdrawn after 6 confirmations",
		"10032": "Please enabled phone/google authenticator",
		"10033": "Fee higher than maximum network transaction fee",
		"10034": "Fee lower than minimum network transaction fee",
		"10035": "Insufficient BTC/LTC",
		"10036": "Withdrawal amount too low",
		"10037": "Trade password not set",
		"10040": "Withdrawal cancellation fails",
		"10041": "Withdrawal address not approved",
		"10042": "Admin password error",
		"10043": "Account equity error, withdrawal failure",
		"10044": "fail to cancel borrowing order",
		"10047": "This function is disabled for sub-account",
		"10100": "User account frozen",
		"10216": "Non-available API",
		"20001": "User does not exist",
		"20002": "Account frozen",
		"20003": "Account frozen due to liquidation",
		"20004": "Futures account frozen",
		"20005": "User futures account does not exist",
		"20006": "Required field missing",
		"20007": "Illegal parameter",
		"20008": "Futures account balance is too low",
		"20009": "Future contract status error",
		"20010": "Risk rate ratio does not exist",
		"20011": `Risk rate higher than 90% before opening position`,
		"20012": `Risk rate higher than 90% after opening position`,
		"20013": "Temporally no counter party price",
		"20014": "System error",
		"20015": "Order does not exist",
		"20016": "Close amount bigger than your open positions",
		"20017": "Not authorized/illegal operation",
		"20018": `Order price differ more than 5% from the price in the last minute`,
		"20019": "IP restricted from accessing the resource",
		"20020": "secretKey does not exist",
		"20021": "Index information does not exist",
		"20022": "Wrong API interface (Cross margin mode shall call cross margin API, fixed margin mode shall call fixed margin API)",
		"20023": "Account in fixed-margin mode",
		"20024": "Signature does not match",
		"20025": "Leverage rate error",
		"20026": "API Permission Error",
		"20027": "No transaction record",
		"20028": "No such contract",
	}
	return errMsgs[code]
}
