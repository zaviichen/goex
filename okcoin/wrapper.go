package okcoin

import (
	"time"
	"github.com/zaviichen/gexch/types"
	"strconv"
)

func TickerFromOKExWs(ex, symbol string, v OKCoinWebsocketFuturesTicker) (types.Ticker) {
	x := types.Ticker{}
	x.Exchange = ex
	x.LocalTime = time.Now()
	x.ServerTime = x.LocalTime
	x.Source = "Ws"
	x.Symbol = symbol
	x.Last = v.Last
	x.Buy = x.Buy
	x.Sell = v.Sell
	x.High = v.High
	x.Low = v.Low
	x.Volume = v.Volume
	return x
}

func OrderBookFromOKExWs(ex, symbol string, v OKCoinWebsocketOrderbook) (types.OrderBook) {
	x := types.OrderBook{}
	x.Exchange = ex
	x.LocalTime = time.Now()
	x.ServerTime = x.LocalTime
	x.Source = "Ws"
	x.Symbol = symbol
	for _, e := range v.Asks {
		ap, _ := strconv.ParseFloat(e[0], 64)
		x.AskPrices = append(x.AskPrices, ap)
		as, _ := strconv.ParseFloat(e[1], 64)
		x.AskSizes = append(x.AskSizes, as)
	}
	for _, e := range v.Bids {
		bp, _ := strconv.ParseFloat(e[0], 64)
		x.BidPrices = append(x.BidPrices, bp)
		bs, _ := strconv.ParseFloat(e[1], 64)
		x.BidSizes = append(x.BidSizes, bs)
	}
	return x
}
