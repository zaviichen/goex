package main

import (
	"fmt"
)

func str(a ...interface{}) string {
	fmt.Println(fmt.Sprint(a))
	return fmt.Sprint(a)
}

func main() {
	//api := New(http.DefaultClient, "", "")
	//
	//ticker, err := api.GetTicker(BTC_CNY)
	//fmt.Println(ticker)
	//fmt.Println(err)

	a := 1
	b := str(a)
	fmt.Println(b)
}

