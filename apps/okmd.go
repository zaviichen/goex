package main

import (
	"fmt"
	"log"
	"github.com/zaviichen/goex/okcoin"
	"github.com/zaviichen/goex/common"
	"time"
	"errors"
	"strings"
)

func WsProcess(o okcoin.OKEx, resp interface{}) (interface{}, error) {
	z := resp.(map[string]interface{})
	channel := z["channel"].(string)
	success := z["success"]
	errorcode := z["errorcode"]
	data, err := common.JSONEncode(z["data"])

	if err != nil {
		return nil, err
	}

	if success != "true" && success != nil {
		errorCodeStr, ok := errorcode.(string)
		if !ok {
			log.Printf("%s Websocket: Unable to convert errorcode to string.", o.Name)
			log.Printf("%s Websocket: channel %s error code: %s.", o.Name, channel, errorcode)
		} else {
			log.Printf("%s Websocket: channel %s error: %s.", o.Name, channel, okcoin.WebsocketErrors[errorCodeStr])
		}
		return nil, errors.New(fmt.Sprintf("channel %s error: %s", channel, okcoin.WebsocketErrors[errorCodeStr]))
	}

	if success == "true" && data == nil {
		return nil, errors.New("no data found")
	}

	log.Printf("%s|%s", channel, string(data))

	switch true {
	case channel == "addChannel":
		return channel, nil
	case strings.Contains(channel, "ok_sub_futureusd"):
		if strings.Contains(channel, "ticker") {
			v := okcoin.OKCoinWebsocketFuturesTicker{}
			err := common.JSONDecode(data, &v)
			if err == nil {
				return okcoin.TickerFromOKExWs(o.Name, "", v), nil
			}
			return nil, err
		} else if strings.Contains(channel, "depth") {
			v := okcoin.OKCoinWebsocketOrderbook{}
			err := common.JSONDecode(data, &v)
			if err == nil {
				return okcoin.OrderBookFromOKExWs(o.Name, "", v), nil
			}
			return nil, err
		} else if strings.Contains(channel, "trade") {
			v := [][]string{}
			err := common.JSONDecode(data, &v)
			return v, err
		} else {
			goto NoChannelTypeMatched
		}
	case strings.Contains(channel, "ok_sub_future") && strings.Contains(channel, "depth"):
		v := okcoin.OKCoinWebsocketOrderbook{}
		err := common.JSONDecode(data, &v)
		return v, err
	default:
		goto NoChannelTypeMatched
	}

NoChannelTypeMatched:
	return channel, errors.New(fmt.Sprintf("unknown channel type: %s", channel))
}

func SubscribeOKExWebSocket() {
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

func main() {
	SubscribeOKExWebSocket()

}
