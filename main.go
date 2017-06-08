package main

import (
	"net/http"
	"fmt"
	. "gexch/okcoin"
)

func main() {
	api := New(http.DefaultClient, "", "")

	ticker, err := api.GetTicker(BTC_CNY)
	fmt.Println(ticker)
	fmt.Println(err)
}

