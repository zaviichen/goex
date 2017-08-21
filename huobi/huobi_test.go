package huobi

import (
	"net/http"
	"testing"
	. "goex/common"
)

var ccy CurrencyPair = BTC_CNY
var ex = NewHuobi(http.DefaultClient, HuobiApiKey, HuobiSecretKey)

func TestHuobi_GetTicker(t *testing.T) {
	dat, err := ex.GetTicker(ccy)
	if err == nil {
		t.Logf("Ticker: %+v", dat)
	} else {
		t.Error(err)
	}
}

func TestHuobi_GetDepth(t *testing.T) {
	dat, err := ex.GetDepth(ccy, 5)
	if err == nil {
		t.Logf("Depth: %+v", dat)
	} else {
		t.Error(err)
	}
}

func TestHuobi_GetTrades(t *testing.T) {
	dat, err := ex.GetTrades(ccy)
	if err == nil {
		t.Logf("Trades: %+v", dat)
	} else {
		t.Error(err)
	}
}

func TestHuobi_GetAccount(t *testing.T) {
	dat, err := ex.GetAccount()
	if err == nil {
		t.Logf("Account: %+v", dat)
	} else {
		t.Error(err)
	}
}
