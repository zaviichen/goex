package okcoin
//
//import (
//	"net/http"
//	"testing"
//	"fmt"
//	"time"
//	. "gexch/common"
//)
//
//var ccy CurrencyPair = BTC_CNY
//var okspot = NewOKCn(http.DefaultClient, OKCnApiKey, OKCnSecretKey)
//
////var ccy CurrencyPair = BTC_USD
////var okspot = NewOKCom(http.DefaultClient, OKComApiKey, OKComSecretKey)
//
//func TestOKSpot_GetTicker(t *testing.T) {
//	ticker, err := okspot.GetTicker(ccy)
//	if err == nil {
//		t.Logf("Ticker: %+v", ticker)
//	} else {
//		t.Error(err)
//	}
//}
//
//func TestOKSpot_GetDepth(t *testing.T) {
//	depth, err := okspot.GetDepth(ccy, 5)
//	if err == nil {
//		t.Logf("Depth: %+v", depth)
//	} else {
//		t.Error(err)
//	}
//}
//
//func TestOKSpot_GetTrades(t *testing.T) {
//	trades, err := okspot.GetTrades(ccy, 10)
//	if err == nil {
//		t.Logf("Trades: %+v", trades)
//	} else {
//		t.Error(err)
//	}
//}
//
//func TestOKSpot_GetAccount(t *testing.T) {
//	account, err := okspot.GetAccount()
//	if err == nil {
//		t.Logf("Account: %+v", account)
//	} else {
//		t.Error(err)
//	}
//}
//
//func TestOKSpot_BuyOrders(t *testing.T) {
//	depth, err := okspot.GetDepth(ccy, 10)
//
//	ob := depth.BidList[len(depth.BidList)-1]
//	order, err := okspot.SendOrder(ccy, "buy", fmt.Sprint(ob.Price), "0.01")
//	if err == nil {
//		t.Logf("Send BuyOrder: %+v", order)
//	} else {
//		t.Error(t)
//	}
//
//	orders, err := okspot.GetOrders(ccy, fmt.Sprint(order.OrderID))
//	if err == nil {
//		t.Logf("Gey BuyOrders: %+v", orders)
//	} else {
//		t.Error(t)
//	}
//
//	openOrders, err := okspot.GetOpenOrders(ccy)
//	if err == nil {
//		t.Logf("Gey OpenOrders: %+v", openOrders)
//	} else {
//		t.Error(t)
//	}
//
//	status, err := okspot.CancelOrder(ccy, fmt.Sprint(order.OrderID))
//	if status == true {
//		t.Logf("Cancel BuyOrder: %+v", status)
//	} else {
//		t.Error(t)
//	}
//
//	time.Sleep(2 * time.Second)
//	openOrders, err = okspot.GetOpenOrders(ccy)
//	if err == nil {
//		t.Logf("Gey OpenOrders after cancellation: %+v", openOrders)
//	} else {
//		t.Error(t)
//	}
//}
//
//func TestOKSpot_SellOrders(t *testing.T) {
//	depth, err := okspot.GetDepth(ccy, 10)
//
//	ob := depth.AskList[len(depth.AskList)-1]
//	order, err := okspot.SendOrder(ccy, "sell", fmt.Sprint(ob.Price), "0.01")
//	if err == nil {
//		t.Logf("Send SelOrder: %+v", order)
//	} else {
//		t.Error(t)
//	}
//
//	orders, err := okspot.GetOrders(ccy, fmt.Sprint(order.OrderID))
//	if err == nil {
//		t.Logf("Gey SelOrders: %+v", orders)
//	} else {
//		t.Error(t)
//	}
//
//	openOrders, err := okspot.GetOpenOrders(ccy)
//	if err == nil {
//		t.Logf("Gey OpenOrders: %+v", openOrders)
//	} else {
//		t.Error(t)
//	}
//
//	status, err := okspot.CancelOrder(ccy, fmt.Sprint(order.OrderID))
//	if status == true {
//		t.Logf("Cancel SelOrder: %+v", status)
//	} else {
//		t.Error(t)
//	}
//
//	time.Sleep(2 * time.Second)
//	openOrders, err = okspot.GetOpenOrders(ccy)
//	if err == nil {
//		t.Logf("Gey OpenOrders after cancellation: %+v", openOrders)
//	} else {
//		t.Error(t)
//	}
//}
//
//func TestOKSpot_GetOrderHistory(t *testing.T) {
//	orders, err := okspot.GetOrderHistory(ccy, 1, 10)
//	if err == nil {
//		t.Logf("HistoryOrders: %+v", orders)
//	} else {
//		t.Error(err)
//	}
//}
//
//func TestOKSpot_GetTradeHistory(t *testing.T) {
//	trades, err := okspot.GetTradeHistory(ccy, 10)
//	if err == nil {
//		t.Logf("HistoryTrades: %+v", trades)
//	} else {
//		t.Error(err)
//	}
//}
