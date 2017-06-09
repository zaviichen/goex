package okcoin

import (
	"net/http"
	"testing"
)

var fccy CurrencyPair = BTC_USD
var fcontract = Weekly
var okex = NewOKEx(http.DefaultClient, OKComApiKey, OKComSecretKey)

func TestOKEx_GetFutTicker(t *testing.T) {
	ticker, err := okex.GetFutTicker(fccy, fcontract)
	if err == nil {
		t.Logf("Ticker: %+v", ticker)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFutDepth(t *testing.T) {
	depth, err := okex.GetFutDepth(fccy, fcontract, 5)
	if err == nil {
		t.Logf("Depth: %+v", depth)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFutTrades(t *testing.T) {
	trades, err := okex.GetFutTrades(fccy, fcontract)
	if err == nil {
		t.Logf("Trades: %+v", trades)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFutIndex(t *testing.T) {
	index, err := okex.GetFutIndex(fccy)
	if err == nil {
		t.Logf("Index: %+v", index)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFutEstimatedPrice(t *testing.T) {
	index, err := okex.GetFutEstimatedPrice(fccy)
	if err == nil {
		t.Logf("Forecast Price: %+v", index)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFutureInfo(t *testing.T) {
	info, err := okex.GetFutureInfo(fccy, fcontract)
	if err == nil {
		t.Logf("Future Info: %+v", info)
	} else {
		t.Error(err)
	}
}