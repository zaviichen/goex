package okcoin

import (
	"testing"
	"github.com/zaviichen/goex/common"
	"time"
)

var ccy string = "btc_cny"
var okspot = NewOKCn(common.OKCnApiKey, common.OKCnSecretKey)

func TestOKCoin_GetTicker(t *testing.T) {
	v, err := okspot.GetTicker(ccy)
	if err == nil {
		t.Logf("GetTicker: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKCoin_GetOrderBook(t *testing.T) {
	v, err := okspot.GetOrderBook(ccy, 5, false)
	if err == nil {
		t.Logf("GetOrderBook: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKCoin_GetTrades(t *testing.T) {
	v, err := okspot.GetTrades(ccy, 0)
	if err == nil {
		t.Logf("GetTrades: %+v", v[:5])
	} else {
		t.Error(err)
	}
}

func TestOKCoin_GetKline(t *testing.T) {
	v, err := okspot.GetKline(ccy, "1min", 5, 0)
	if err == nil {
		t.Logf("GetKline: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKCoin_GetUserInfo(t *testing.T) {
	v, err := okspot.GetUserInfo()
	if err == nil {
		t.Logf("GetUserInfo: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKCoin_PlaceOrders(t *testing.T) {
	depth, err := okspot.GetOrderBook(ccy, 10, false)

	ob := depth.Bids[len(depth.Bids)-1]
	orderid, err := okspot.Trade(ccy, "buy", ob[0], 0.01)
	if err == nil {
		t.Logf("Send BuyOrder: %+v", orderid)
	} else {
		t.Error(err)
	}

	time.Sleep(200 * time.Millisecond)
	orders, err := okspot.GetOrderInfo(ccy, orderid)
	if err == nil {
		t.Logf("Get BuyOrder: %+v", orders)
	} else {
		t.Error(err)
	}

	orderid, err = okspot.CancelOrder(ccy, orderid)
	if err == nil {
		t.Logf("Cancel BuyOrder: %+v", orderid)
	} else {
		t.Error(err)
	}
}

func TestOKCoin_GetOrderHistory(t *testing.T) {
	v, err := okspot.GetOrderHistory(1,1,"1", ccy)
	if err == nil {
		t.Logf("GetOrderHistory: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKCoin_GetTradeHistory(t *testing.T) {
	v, err := okspot.GetTradeHistory(ccy, 7564501566)
	if err == nil {
		t.Logf("GetTradeHistory: %+v", v[:5])
	} else {
		t.Error(err)
	}
}
