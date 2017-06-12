package okcoin

import (
	"net/http"
	"testing"
	"time"
	. "gexch/common"
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

func TestOKEx_GetFutAccount(t *testing.T) {
	dat, err := okex.GetFutAccount()
	if err == nil {
		t.Logf("Future Account: %+v", dat)
	} else {
		t.Error(err)
	}
}

func TestOKEx_GetFutPosition(t *testing.T) {
	dat, err := okex.GetFutPosition(fccy, fcontract)
	if err == nil {
		t.Logf("Future Position: %+v", dat)
	} else {
		t.Error(err)
	}
}

func TestOKEx_BuyOrders(t *testing.T) {
	depth, err := okex.GetFutDepth(fccy, fcontract, 10)

	ob := depth.BidList[len(depth.BidList)-1]
	order, err := okex.SendFutOrder(fccy, fcontract, ob.Price, 1, 1, 0, 20)
	if err == nil {
		t.Logf("Send BuyOrder: %+v", order)
	} else {
		t.Error(err)
	}

	ids := []string{order.OrderID}
	orders, err := okex.GetFutOrders(fccy, fcontract, ids)
	if err == nil {
		t.Logf("Gey BuyOrders: %+v", orders)
	} else {
		t.Error(err)
	}

	openOrders, err := okex.GetFutOpenOrders(fccy, fcontract)
	if err == nil {
		t.Logf("Gey OpenOrders: %+v", openOrders)
	} else {
		t.Error(err)
	}

	status, err := okex.CancelFutOrder(fccy, fcontract, order.OrderID)
	if status == true {
		t.Logf("Cancel BuyOrder: %+v", status)
	} else {
		t.Error(err)
	}

	time.Sleep(2 * time.Second)
	openOrders, err = okex.GetFutOpenOrders(fccy, fcontract)
	if err == nil {
		t.Logf("Gey OpenOrders after cancellation: %+v", openOrders)
	} else {
		t.Error(err)
	}
}

func TestOKEx_SellOrders(t *testing.T) {
	depth, err := okex.GetFutDepth(fccy, fcontract, 10)

	ob := depth.AskList[len(depth.AskList)-1]
	order, err := okex.SendFutOrder(fccy, fcontract, ob.Price, 1, 2, 0, 20)
	if err == nil {
		t.Logf("Send SellOrder: %+v", order)
	} else {
		t.Error(err)
	}

	ids := []string{order.OrderID}
	orders, err := okex.GetFutOrders(fccy, fcontract, ids)
	if err == nil {
		t.Logf("Gey BuyOrders: %+v", orders)
	} else {
		t.Error(err)
	}

	openOrders, err := okex.GetFutOpenOrders(fccy, fcontract)
	if err == nil {
		t.Logf("Gey OpenOrders: %+v", openOrders)
	} else {
		t.Error(err)
	}

	status, err := okex.CancelFutOrder(fccy, fcontract, order.OrderID)
	if status == true {
		t.Logf("Cancel BuyOrder: %+v", status)
	} else {
		t.Error(err)
	}

	time.Sleep(2 * time.Second)
	openOrders, err = okex.GetFutOpenOrders(fccy, fcontract)
	if err == nil {
		t.Logf("Gey OpenOrders after cancellation: %+v", openOrders)
	} else {
		t.Error(err)
	}
}
