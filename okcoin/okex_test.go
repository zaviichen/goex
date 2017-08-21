package okcoin

import (
	"testing"
	"github.com/zaviichen/goex/common"
	"time"
)

var fccy string = "btc_usd"
var fcontract = Weekly
var okex = NewOKEx(common.OKComApiKey, common.OKComSecretKey)

func TestOKEx_GetFuturesTicker(t *testing.T) {
	v, err := okex.GetFuturesTicker(fccy, fcontract)
	if err == nil {
		t.Logf("GetFuturesTicker: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesDepth(t *testing.T) {
	v, err := okex.GetFuturesDepth(fccy, fcontract, 5, false)
	if err == nil {
		t.Logf("GetFuturesDepth: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesTrades(t *testing.T) {
	v, err := okex.GetFuturesTrades(fccy, fcontract)
	if err == nil {
		t.Logf("GetFuturesTrades: #=%d, %+v", len(v), v[:5])
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesIndex(t *testing.T) {
	v, err := okex.GetFuturesIndex(fccy)
	if err == nil {
		t.Logf("GetFuturesIndex: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesExchangeRate(t *testing.T) {
	v, err := okex.GetFuturesExchangeRate()
	if err == nil {
		t.Logf("GetFuturesExchangeRate: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesEstimatedPrice(t *testing.T) {
	v, err := okex.GetFuturesEstimatedPrice(fccy)
	if err == nil {
		t.Logf("GetFuturesEstimatedPrice: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesKline(t *testing.T) {
	v, err := okex.GetFuturesKline(fccy, "1min", fcontract, 5, 0)
	if err == nil {
		t.Logf("GetFuturesKline: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesHoldAmount(t *testing.T) {
	v, err := okex.GetFuturesHoldAmount(fccy, fcontract)
	if err == nil {
		t.Logf("GetFuturesHoldAmount: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesUserInfo(t *testing.T) {
	v, err := okex.GetFuturesUserInfo()
	if err == nil {
		t.Logf("GetFuturesUserInfo: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFuturesPosition(t *testing.T) {
	v, err := okex.GetFuturesPosition(fccy, fcontract)
	if err == nil {
		t.Logf("GetFuturesPosition: %+v", v)
	} else {
		t.Error(err)
	}
}

func TestOKEx_PlaceOrders(t *testing.T) {
	depth, err := okex.GetFuturesDepth(fccy, fcontract, 10, false)

	ob := depth.Bids[len(depth.Bids)-1]
	orderid, err := okex.FuturesTrade(fccy, fcontract, OpenLong, ob[0], 1, 0, 20)
	if err == nil {
		t.Logf("Send BuyOrder: %+v", orderid)
	} else {
		t.Error(err)
	}

	time.Sleep(200 * time.Millisecond)
	orders, err := okex.GetFuturesOrdersInfo(fccy, fcontract, orderid)
	if err == nil {
		t.Logf("Get BuyOrder: %+v", orders)
	} else {
		t.Error(err)
	}

	orderid, err = okex.CancelFuturesOrder(fccy, fcontract, orders[0].OrderID)
	if err == nil {
		t.Logf("Cancel BuyOrder: %+v", orderid)
	} else {
		t.Error(err)
	}
}
