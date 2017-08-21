package huobi

import (
	"net/http"
	"net/url"
	"fmt"
	"time"
	"errors"
	"strconv"
	. "goex/common"
	"strings"
	"encoding/json"
)

const (
	HuobiName     = "huobi.com"
	BaseUri       = "https://api.huobi.com/"
	TradeApiV3Uri = BaseUri + "apiv3/"
	TickerUri     = "staticmarket/ticker_%s_json.js"
	DepthUri      = "staticmarket/depth_%s_%d.js"
	TradesUri     = "staticmarket/detail_%s_json.js"
)

type Huobi struct {
	ExchangeBase
}

func NewHuobi(client *http.Client, api string, secret string) *Huobi {
	ex := new(Huobi)
	ex.Name = HuobiName
	ex.BaseUri = BaseUri
	ex.Enable = true
	ex.HttpClient = client
	ex.ApiKey = api
	ex.SecretKey = secret
	ex.MakerFee = 0.004
	ex.TakerFee = 0.004
	return ex
}

func (ex *Huobi) GetTicker(ccy CurrencyPair) (*Ticker, error) {
	url := ex.BaseUri + TickerUri
	switch ccy {
	case BTC_CNY:
		url = fmt.Sprintf(url, "btc")
	case LTC_CNY:
		url = fmt.Sprintf(url, "ltc")
	default:
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	dat, err := HttpGet(ex.HttpClient, url)
	if err != nil {
		return nil, err
	}

	tmap := dat["ticker"].(map[string]interface{})
	ticker := new(Ticker)
	ticker.Date, _ = strconv.ParseUint(dat["time"].(string), 10, 64)
	ticker.Last = tmap["last"].(float64)
	ticker.Buy = tmap["buy"].(float64)
	ticker.Sell = tmap["sell"].(float64)
	ticker.Low = tmap["low"].(float64)
	ticker.High = tmap["high"].(float64)
	ticker.Vol = tmap["vol"].(float64)
	return ticker, nil
}

func (ex *Huobi) GetDepth(ccy CurrencyPair, size int) (*Depth, error) {
	url := ex.BaseUri + DepthUri
	switch ccy {
	case BTC_CNY:
		url = fmt.Sprintf(url, "btc", size)
	case LTC_CNY:
		url = fmt.Sprintf(url, "ltc", size)
	default:
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	dat, err := HttpGet(ex.HttpClient, url)
	if err != nil {
		return nil, err
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

func (ex *Huobi) GetTrades(ccy CurrencyPair) ([]Trade, error) {
	url := ex.BaseUri + TradesUri
	switch ccy {
	case BTC_CNY:
		url = fmt.Sprintf(url, "btc")
	case LTC_CNY:
		url = fmt.Sprintf(url, "ltc")
	default:
		return nil, errors.New("Unsupport The CurrencyPair")
	}

	dat, err := HttpGet(ex.HttpClient, url)
	if err != nil {
		return nil, err
	}

	tradesmap := dat["trades"].([]interface{})
	now := time.Now()

	var trades []Trade
	for _, t := range tradesmap {
		tr := t.(map[string]interface{})
		trade := Trade{}
		trade.Amount = tr["amount"].(float64)
		trade.Price = tr["price"].(float64)
		trade.Type = tr["type"].(string)
		timeStr := tr["time"].(string)
		timeMeta := strings.Split(timeStr, ":")
		h, _ := strconv.Atoi(timeMeta[0])
		m, _ := strconv.Atoi(timeMeta[1])
		s, _ := strconv.Atoi(timeMeta[2])
		if now.Hour() == 0 {
			if h <= 23 && h >= 20 {
				pre := now.AddDate(0, 0, -1)
				trade.Date = time.Date(pre.Year(), pre.Month(), pre.Day(), h, m, s, 0, time.Local).Unix() * 1000
			} else if h == 0 {
				trade.Date = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local).Unix() * 1000
			}
		} else {
			trade.Date = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local).Unix() * 1000
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

func (ex *Huobi) PostRequest(method string, v url.Values, result interface{}) (err error) {
	v.Set("method", method)
	v.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	v.Set("access_key", ex.ApiKey)

	hash := GetMD5([]byte(v.Encode() + "&secret_key=" + ex.SecretKey))
	v.Set("sign", StringToLower(HexEncodeToString(hash)))

	body, err := HttpPostForm(ex.HttpClient, TradeApiV3Uri, v)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		println(string(body))
		return err
	}
	return nil
}

func (ex *Huobi) GetAccount() (*Account, error) {
	v := url.Values{}

	x := new(_AccountInfo)
	err := ex.PostRequest("get_account_info", v, x)
	if err != nil {
		return nil, err
	}

	acc := new(Account)
	acc.Exchange = ex.Name
	acc.Asset = x.Total
	acc.NetAsset = x.NetAsset
	acc.SubAccounts = make(map[Currency]SubAccount, 3)
	acc.SubAccounts[BTC] = SubAccount{BTC, x.AvailableBtc, x.FrozenBtc, x.LoanBtc}
	acc.SubAccounts[LTC] = SubAccount{LTC, x.AvailableLtc, x.FrozenLtc, x.LoanLtc}
	acc.SubAccounts[CNY] = SubAccount{CNY, x.AvailableCny, x.FrozenCny, x.LoanCny}
	return acc, nil
}
