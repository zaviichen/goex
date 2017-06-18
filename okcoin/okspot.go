package okcoin
//
//import (
//	"net/http"
//	"net/url"
//	"strings"
//	"strconv"
//	"fmt"
//	"encoding/json"
//	"errors"
//	. "gexch/common"
//)
//
//const (
//	OKCnName         = "OKCn"
//	OKComName        = "OKCom"
//	OKCnBaseUri      = "https://www.okcoin.cn/api/v1/"
//	OKComBaseUri     = "https://www.okcoin.com/api/v1/"
//	TickerUri        = "ticker.do?symbol=%s"
//	DepthUri         = "depth.do?symbol=%s&size=%d"
//	TradesUri        = "trades.do?symbol=%s&since=%d"
//	UserInfoUri      = "userinfo.do"
//	TradeUri         = "trade.do"
//	CancelOrderUri   = "cancel_order.do"
//	OrderInfoUri     = "order_info.do"
//	OrderHistoryUri  = "order_history.do"
//	TradeHistoryUri = "trade_history.do"
//)
//
//var ExPairSymbol = map[CurrencyPair]string{
//	BTC_CNY: "btc_cny",
//	BTC_USD: "btc_usd",
//	LTC_CNY: "ltc_cny",
//	LTC_USD: "ltc_usd",
//	ETH_CNY: "eth_cny",
//	ETH_USD: "eth_usd",
//	ETH_BTC: "eth_btc",
//	ETC_CNY: "etc_cny",
//	ETC_USD: "etc_usd",
//	ETC_BTC: "etc_btc"}
//
//type OKSpot struct {
//	client    *http.Client
//	apiKey    string
//	secretKey string
//	Name      string
//	BaseUri   string
//}
//
//func NewOKCn(client *http.Client, api string, secret string) *OKSpot {
//	ex := new(OKSpot)
//	ex.client = client
//	ex.apiKey = api
//	ex.secretKey = secret
//	ex.Name = OKCnName
//	ex.BaseUri = OKCnBaseUri
//	return ex
//}
//
//func NewOKCom(client *http.Client, api string, secret string) *OKSpot {
//	ex := new(OKSpot)
//	ex.client = client
//	ex.apiKey = api
//	ex.secretKey = secret
//	ex.Name = OKComName
//	ex.BaseUri = OKComBaseUri
//	return ex
//}
//
////func GetParamMD5Sign(_ string, params string) (string, error) {
////	hash := md5.New()
////	_, err := hash.Write([]byte(params))
////	if err != nil {
////		return "", err
////	}
////	return hex.EncodeToString(hash.Sum(nil)), nil
////}
//
//func BuildPostForm(postForm *url.Values, apiKey, secretKey string) error {
//	postForm.Set("api_key", apiKey)
//
//	payload := postForm.Encode()
//	payload = payload + "&secret_key=" + secretKey
//
//	sign, err := GetParamMD5Sign(secretKey, payload)
//	if err != nil {
//		return err
//	}
//
//	postForm.Set("sign", strings.ToUpper(sign))
//	return nil
//}
//
//func (ctx *OKSpot) GetTicker(currency CurrencyPair) (*Ticker, error) {
//	url := fmt.Sprintf(ctx.BaseUri+TickerUri, ExPairSymbol[currency])
//	dat, err := HttpGet(ctx.client, url)
//	if err != nil {
//		return nil, err
//	}
//
//	tmap := dat["ticker"].(map[string]interface{})
//	ticker := new(Ticker)
//	ticker.Date, _ = strconv.ParseUint(dat["date"].(string), 10, 64)
//	ticker.Last, _ = strconv.ParseFloat(tmap["last"].(string), 64)
//	ticker.Buy, _ = strconv.ParseFloat(tmap["buy"].(string), 64)
//	ticker.Sell, _ = strconv.ParseFloat(tmap["sell"].(string), 64)
//	ticker.Low, _ = strconv.ParseFloat(tmap["low"].(string), 64)
//	ticker.High, _ = strconv.ParseFloat(tmap["high"].(string), 64)
//	ticker.Vol, _ = strconv.ParseFloat(tmap["vol"].(string), 64)
//	return ticker, nil
//}
//
//func (ctx *OKSpot) GetDepth(currency CurrencyPair, size int) (*Depth, error) {
//	url := fmt.Sprintf(ctx.BaseUri+DepthUri, ExPairSymbol[currency], size)
//	dat, err := HttpGet(ctx.client, url)
//	if err != nil {
//		return nil, err
//	}
//
//	if dat["result"] != nil && !dat["result"].(bool) {
//		return nil, errors.New(fmt.Sprintf("%.0f", dat["error_code"].(float64)))
//	}
//
//	var depth Depth
//	for _, v := range dat["asks"].([]interface{}) {
//		var dr DepthRecord
//		for i, vv := range v.([]interface{}) {
//			switch i {
//			case 0:
//				dr.Price = vv.(float64)
//			case 1:
//				dr.Amount = vv.(float64)
//			}
//		}
//		depth.AskList = append(depth.AskList, dr)
//	}
//
//	for _, v := range dat["bids"].([]interface{}) {
//		var dr DepthRecord
//		for i, vv := range v.([]interface{}) {
//			switch i {
//			case 0:
//				dr.Price = vv.(float64)
//			case 1:
//				dr.Amount = vv.(float64)
//			}
//		}
//		depth.BidList = append(depth.BidList, dr)
//	}
//
//	return &depth, nil
//}
//
//func (ctx *OKSpot) GetTrades(currency CurrencyPair, since int64) ([]Trade, error) {
//	url := fmt.Sprintf(ctx.BaseUri+TradesUri, ExPairSymbol[currency], since)
//	dat, err := HttpGet2(ctx.client, url)
//	if err != nil {
//		return nil, err
//	}
//
//	var trades []Trade
//	err = json.Unmarshal(dat, &trades)
//	if err != nil {
//		return nil, err
//	}
//
//	return trades, nil
//}
//
//func (ctx *OKSpot) GetAccount() (*Account, error) {
//	postData := url.Values{}
//	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
//	if err != nil {
//		return nil, err
//	}
//
//	body, err := HttpPostForm(ctx.client, ctx.BaseUri+UserInfoUri, postData)
//	if err != nil {
//		return nil, err
//	}
//
//	var respMap map[string]interface{}
//	err = json.Unmarshal(body, &respMap)
//	if err != nil {
//		return nil, err
//	}
//
//	if !respMap["result"].(bool) {
//		errCode := strconv.FormatFloat(respMap["error_code"].(float64), 'f', 0, 64)
//		return nil, errors.New(errCode)
//	}
//
//	info := respMap["info"].(map[string]interface{})
//	funds := info["funds"].(map[string]interface{})
//	asset := funds["asset"].(map[string]interface{})
//	free := funds["free"].(map[string]interface{})
//	freezed := funds["freezed"].(map[string]interface{})
//
//	account := new(Account)
//	account.Exchange = ctx.Name
//	account.Asset, _ = strconv.ParseFloat(asset["total"].(string), 64)
//	account.NetAsset, _ = strconv.ParseFloat(asset["net"].(string), 64)
//
//	var btc SubAccount
//	var ltc SubAccount
//	var eth SubAccount
//
//	btc.Currency = BTC
//	btc.Amount, _ = strconv.ParseFloat(free["btc"].(string), 64)
//	btc.LoanAmount = 0
//	btc.ForzenAmount, _ = strconv.ParseFloat(freezed["btc"].(string), 64)
//
//	ltc.Currency = LTC
//	ltc.Amount, _ = strconv.ParseFloat(free["ltc"].(string), 64)
//	ltc.LoanAmount = 0
//	ltc.ForzenAmount, _ = strconv.ParseFloat(freezed["ltc"].(string), 64)
//
//	eth.Currency = ETH
//	eth.Amount, _ = strconv.ParseFloat(free["ltc"].(string), 64)
//	eth.LoanAmount = 0
//	eth.ForzenAmount, _ = strconv.ParseFloat(freezed["ltc"].(string), 64)
//
//	account.SubAccounts = make(map[Currency]SubAccount, 3)
//	account.SubAccounts[BTC] = btc
//	account.SubAccounts[LTC] = ltc
//	account.SubAccounts[ETH] = eth
//
//	var fiat SubAccount
//	if ctx.Name == OKCnName {
//		fiat.Currency = CNY
//		fiat.Amount, _ = strconv.ParseFloat(free["cny"].(string), 64)
//		fiat.LoanAmount = 0
//		fiat.ForzenAmount, _ = strconv.ParseFloat(freezed["cny"].(string), 64)
//		account.SubAccounts[CNY] = fiat
//	} else if ctx.Name == OKComName {
//		fiat.Currency = USD
//		fiat.Amount, _ = strconv.ParseFloat(free["usd"].(string), 64)
//		fiat.LoanAmount = 0
//		fiat.ForzenAmount, _ = strconv.ParseFloat(freezed["usd"].(string), 64)
//		account.SubAccounts[USD] = fiat
//	}
//	return account, nil
//}
//
//func (ctx *OKSpot) SendOrder(currency CurrencyPair, side, price, amount string) (*Order, error) {
//	postData := url.Values{}
//	postData.Set("type", side)
//
//	if side != "buy_market" {
//		postData.Set("amount", amount)
//	}
//	if side != "sell_market" {
//		postData.Set("price", price)
//	}
//	postData.Set("symbol", ExPairSymbol[currency])
//
//	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
//	if err != nil {
//		return nil, err
//	}
//
//	body, err := HttpPostForm(ctx.client, ctx.BaseUri+TradeUri, postData)
//	if err != nil {
//		return nil, err
//	}
//
//	var respMap map[string]interface{}
//	err = json.Unmarshal(body, &respMap)
//	if err != nil {
//		return nil, err
//	}
//
//	if !respMap["result"].(bool) {
//		return nil, errors.New(string(body))
//	}
//
//	order := new(Order)
//	order.OrderID = int(respMap["order_id"].(float64))
//	order.Price, _ = strconv.ParseFloat(price, 64)
//	order.Amount, _ = strconv.ParseFloat(amount, 64)
//	order.Currency = currency
//	order.Status = ORDER_UNFINISH
//
//	switch side {
//	case "buy":
//		order.Side = BUY
//	case "sell":
//		order.Side = SELL
//	}
//
//	return order, nil
//}
//
//func (ctx *OKSpot) CancelOrder(currency CurrencyPair, orderId string) (bool, error) {
//	postData := url.Values{}
//	postData.Set("order_id", orderId)
//	postData.Set("symbol", ExPairSymbol[currency])
//
//	BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
//	body, err := HttpPostForm(ctx.client, ctx.BaseUri+CancelOrderUri, postData)
//	if err != nil {
//		return false, err
//	}
//
//	var respMap map[string]interface{}
//	err = json.Unmarshal(body, &respMap)
//	if err != nil {
//		return false, err
//	}
//
//	if !respMap["result"].(bool) {
//		return false, errors.New(string(body))
//	}
//	return true, nil
//}
//
//func fillOrder(currency CurrencyPair, orderMap map[string]interface{}) Order {
//	var order Order
//	order.Currency = currency
//	order.OrderID = int(orderMap["order_id"].(float64))
//	order.Amount = orderMap["amount"].(float64)
//	order.Price = orderMap["price"].(float64)
//	order.DealAmount = orderMap["deal_amount"].(float64)
//	order.AvgPrice = orderMap["avg_price"].(float64)
//	order.OrderTime = int(orderMap["create_date"].(float64))
//
//	//status:-1:已撤销  0:未成交  1:部分成交  2:完全成交 4:撤单处理中
//	switch int(orderMap["status"].(float64)) {
//	case -1:
//		order.Status = ORDER_CANCEL
//	case 0:
//		order.Status = ORDER_UNFINISH
//	case 1:
//		order.Status = ORDER_PART_FINISH
//	case 2:
//		order.Status = ORDER_FINISH
//	case 4:
//		order.Status = ORDER_CANCEL_ING
//	}
//
//	switch orderMap["type"].(string) {
//	case "buy":
//		order.Side = BUY
//	case "sell":
//		order.Side = SELL
//	case "buy_market":
//		order.Side = BUY_MARKET
//	case "sell_market":
//		order.Side = SELL_MARKET
//	}
//	return order
//}
//
//func (ctx *OKSpot) GetOrders(currency CurrencyPair, orderId string) ([]Order, error) {
//	postData := url.Values{}
//	postData.Set("order_id", orderId)
//	postData.Set("symbol", ExPairSymbol[currency])
//
//	BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
//	body, err := HttpPostForm(ctx.client, ctx.BaseUri+OrderInfoUri, postData)
//	if err != nil {
//		return nil, err
//	}
//
//	var respMap map[string]interface{}
//	err = json.Unmarshal(body, &respMap)
//	if err != nil {
//		return nil, err
//	}
//
//	if !respMap["result"].(bool) {
//		return nil, errors.New(string(body))
//	}
//
//	var orders []Order
//	for _, v := range respMap["orders"].([]interface{}) {
//		order := fillOrder(currency, v.(map[string]interface{}))
//		orders = append(orders, order)
//	}
//	return orders, nil
//}
//
//func (ctx *OKSpot) GetOpenOrders(currency CurrencyPair) ([]Order, error) {
//	return ctx.GetOrders(currency, "-1")
//}
//
//func (ctx *OKSpot) GetOrderHistory(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
//
//	postData := url.Values{}
//	postData.Set("status", "1")
//	postData.Set("symbol", ExPairSymbol[currency])
//	postData.Set("current_page", strconv.Itoa(currentPage))
//	postData.Set("page_length", strconv.Itoa(pageSize))
//
//	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
//	body, err := HttpPostForm(ctx.client, ctx.BaseUri+OrderHistoryUri, postData)
//	if err != nil {
//		return nil, err
//	}
//
//	var respMap map[string]interface{}
//	err = json.Unmarshal(body, &respMap)
//	if err != nil {
//		return nil, err
//	}
//
//	if !respMap["result"].(bool) {
//		return nil, errors.New(string(body))
//	}
//
//	var orders []Order
//	for _, v := range respMap["orders"].([]interface{}) {
//		order := fillOrder(currency, v.(map[string]interface{}))
//		orders = append(orders, order)
//	}
//	return orders, nil
//}
//
//func (ctx *OKSpot) GetTradeHistory(currency CurrencyPair, since int64) ([]Trade, error) {
//	postData := url.Values{}
//	postData.Set("symbol", ExPairSymbol[currency])
//	postData.Set("since", fmt.Sprintf("%d", since))
//
//	err := BuildPostForm(&postData, ctx.apiKey, ctx.secretKey)
//	body, err := HttpPostForm(ctx.client, ctx.BaseUri+TradeHistoryUri, postData)
//	if err != nil {
//		return nil, err
//	}
//
//	var trades []Trade
//	err = json.Unmarshal(body, &trades)
//	if err != nil {
//		return nil, err
//	}
//	return trades, nil
//}
