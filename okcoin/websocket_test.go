package okcoin

import (
	"testing"
	"github.com/zaviichen/gexch/common"
	//"log"
	"fmt"
	//"reflect"
	//"strconv"
	//"github.com/gorilla/websocket"
	"log"
	"time"
)

func TestOKEx_WsOrders(t *testing.T) {
	apiKey := "f8d3e594-991d-48d5-9560-9d4193733290"
	secretKey := "6595A5938CDB28DA31DA9D39F71FEE7B"
	o := NewOKEx(apiKey, secretKey)

	err := o.WsConnect()

	if err != nil {
		t.Error(err)
		return
	}

	o.AddOrderInfo()
	o.AddUserInfo()

	channels := []string{
		fmt.Sprintf(OKExWebSocketTicker, "ltc", "this_week"),
	}

	for _, channel := range channels {
		o.AddChannel(channel)
	}


	listen := func() {
		for {
			resp := []interface{}{}
			o.WebsocketConn.ReadJSON(&resp)
			fmt.Println(resp)
		}
	}
	go listen()

	time.Sleep(600 * time.Second)
}

func _TestOKEx_SubscribeMarketData(t *testing.T) {
	o := NewOKEx(common.OKComApiKey, common.OKComSecretKey)
	symbol := "btc"
	contract := Weekly

	err := o.WsConnect()

	if err != nil {
		t.Error(err)
		return
	}

	channels := []string{
		fmt.Sprintf(OKExWebSocketTicker, symbol, contract),
		fmt.Sprintf(OKExWebSocketDepthFull, symbol, contract, 10),
		fmt.Sprintf(OKExWebSocketTrade, symbol, contract),
	}

	for _, channel := range channels {
		o.AddChannel(channel)
	}

	listen := func() {
		for {
			resp := []interface{}{}
			o.WebsocketConn.ReadJSON(&resp)
			fmt.Println("a")
			fmt.Println(resp)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, y := range resp {
				v, err := o.WsProcessResponse(y)
				if err == nil {
					log.Printf("WsResponse: %+v", v)
					t.Logf("WsResponse: %+v", v)
				} else {
					t.Error(err)
				}
			}
		}
	}
	go listen()

	time.Sleep(5 * time.Second)
	for _, channel := range channels {
		o.RemoveChannel(channel)
	}
}
