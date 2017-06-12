package huobi

import (
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"fmt"
	"time"
	"errors"
	"strconv"
	. "gexch/common"
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
	client    *http.Client
	apiKey    string
	secretKey string
	Name      string
	BaseUri   string
}

func NewHuobi(client *http.Client, api string, secret string) *Huobi {
	ex := new(Huobi)
	ex.client = client
	ex.apiKey = api
	ex.secretKey = secret
	ex.Name = HuobiName
	ex.BaseUri = BaseUri
	return ex
}

func GetParamMD5Sign(_ string, params string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func BuildPostForm(postForm *url.Values, apiKey, secretKey string) error {
	postForm.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	postForm.Set("access_key", apiKey)
	postForm.Set("secret_key", secretKey)
	sign, err := GetParamMD5Sign(secretKey, postForm.Encode())
	if err != nil {
		return err
	}
	postForm.Set("sign", sign)
	postForm.Del("secret_key")
	return nil
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

	dat, err := HttpGet(ex.client, url)
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

	dat, err := HttpGet(ex.client, url)
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

	dat, err := HttpGet(ex.client, url)
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

func (ex *Huobi) GetAccount() (*Account, error) {
	postData := url.Values{}
	postData.Set("method", "get_account_info")
	postData.Set("created", fmt.Sprintf("%d", time.Now().Unix()))
	postData.Set("access_key", ex.apiKey)
	postData.Set("secret_key", ex.secretKey)

	sign, _ := GetParamMD5Sign(ex.secretKey, postData.Encode())
	postData.Set("sign", sign)
	postData.Del("secret_key")

	body, err := HttpPostForm(ex.client, TradeApiV3Uri, postData)
	if err != nil {
		return nil, err
	}

	var dat map[string]interface{}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		println(string(body))
		return nil, err
	}

	if dat["code"] != nil {
		return nil, errors.New(fmt.Sprintf("%s", dat))
	}

	account := new(Account)
	account.Exchange = ex.Name
	account.Asset, _ = strconv.ParseFloat(dat["total"].(string), 64)
	account.NetAsset, _ = strconv.ParseFloat(dat["net_asset"].(string), 64)
	
	var btc SubAccount
	var ltc SubAccount
	var cny SubAccount
	
	btc.Currency = BTC
	btc.Amount, _ = strconv.ParseFloat(dat["available_btc_display"].(string), 64)
	btc.LoanAmount, _ = strconv.ParseFloat(dat["loan_btc_display"].(string), 64)
	btc.ForzenAmount, _ = strconv.ParseFloat(dat["frozen_btc_display"].(string), 64)
	
	ltc.Currency = LTC
	ltc.Amount, _ = strconv.ParseFloat(dat["available_ltc_display"].(string), 64)
	ltc.LoanAmount, _ = strconv.ParseFloat(dat["loan_ltc_display"].(string), 64)
	ltc.ForzenAmount, _ = strconv.ParseFloat(dat["frozen_ltc_display"].(string), 64)
	
	cny.Currency = CNY
	cny.Amount, _ = strconv.ParseFloat(dat["available_cny_display"].(string), 64)
	cny.LoanAmount, _ = strconv.ParseFloat(dat["loan_cny_display"].(string), 64)
	cny.ForzenAmount, _ = strconv.ParseFloat(dat["frozen_cny_display"].(string), 64)

	account.SubAccounts = make(map[Currency]SubAccount, 3)
	account.SubAccounts[BTC] = btc
	account.SubAccounts[LTC] = ltc
	account.SubAccounts[CNY] = cny
	return account, nil
}
