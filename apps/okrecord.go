package main

import (
	"log"
	"fmt"
	"time"
	"github.com/zaviichen/gexch/okcoin"
	"github.com/zaviichen/gexch/common"
)

func SubscribeOKExWebSocket2() {
	o := okcoin.NewOKEx(common.OKComApiKey, common.OKComSecretKey)
	symbol := "btc"
	contract := okcoin.Weekly

	err := o.WsConnect()
	if err == nil {
		log.Printf("%s Connected to Websocket", o.Name)
	} else {
		log.Printf("%s Unable to connect to Websocket. Error: %s", o.Name, err)
	}

	channels := []string{
		fmt.Sprintf(okcoin.OKExWebSocketTicker, symbol, contract),
		fmt.Sprintf(okcoin.OKExWebSocketDepthFull, symbol, contract, 10),
		fmt.Sprintf(okcoin.OKExWebSocketTrade, symbol, contract),
	}

	for _, channel := range channels {
		o.AddChannel(channel)
	}

	listen := func() {
		for {
			resp := []interface{}{}
			o.WebsocketConn.ReadJSON(&resp)
			fmt.Println(resp)

			if err != nil {
				log.Println(err)
				continue
			}

			for _, y := range resp {
				v, err := WsProcess(*o, y)
				if err == nil {
					log.Printf("WsResponse: %+v", v)
				} else {
					log.Printf("Error: %s", err)
				}
			}
		}
	}
	go listen()

	time.Sleep(600 * time.Second)
	//for _, channel := range channels {
	//	o.RemoveChannel(channel)
	//}
}

