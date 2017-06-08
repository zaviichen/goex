package okcoin

import (
	"net/http"
	"net/url"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

const (
	BaseUri = "https://www.okcoin.cn/api/v1/"
	TickerUri       = "ticker.do"
	DepthUri        = "depth.do"
	TradesUri       = "trades.do"
	url_userinfo      = "userinfo.do"
	url_trade         = "trade.do"
	url_cancel_order  = "cancel_order.do"
	url_order_info    = "order_info.do"
	url_orders_info   = "orders_info.do"
	order_history_uri = "order_history.do"
	trade_uri 	  = "trade_history.do"
)

type OKCoinCn struct {
	client *http.Client
	apiKey string
	secretKey string
}



func New(client *http.Client, api string, secret string) *OKCoinCn {
	return &OKCoinCn{client, api, secret}
}

func (ctx *OKCoinCn) buildPostForm(postForm *url.Values) error {
	postForm.Set("api_key", ctx.apiKey)

	payload := postForm.Encode()
	payload = payload + "&secret_key=" + ctx.secretKey

	sign, err := GetParamMD5Sign(ctx.secretKey, payload)
	if err != nil {
		return err
	}

	postForm.Set("sign", strings.ToUpper(sign))
	return nil
}

func exCoinName(coin CurrencyPair) string {
	switch coin {
	case BTC_CNY:
		return "btc_cny"
	case LTC_CNY:
		return "ltc_cny"
	case BTC_USD:
		return "btc_usd"
	case LTC_USD:
		return "ltc_usd"
	default:
		return ""
	}
}

func GetParamMD5Sign(secret string, params string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (ctx *OKCoinCn) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var tickerMap map[string]interface{}
	var ticker Ticker

	url := BaseUri + TickerUri + "?symbol=" + exCoinName(currency)
	dat, err := HttpGet(ctx.client, url)
	if err != nil {
		return nil, err
	}

	tickerMap = dat["ticker"].(map[string]interface{})
	ticker.Date, _ = strconv.ParseUint(dat["date"].(string), 10, 64)
	ticker.Last, _ = strconv.ParseFloat(tickerMap["last"].(string), 64)
	ticker.Buy, _ = strconv.ParseFloat(tickerMap["buy"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(tickerMap["sell"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(tickerMap["low"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(tickerMap["high"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(tickerMap["vol"].(string), 64)

	return &ticker, nil
}
