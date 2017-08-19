package okcoin

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zaviichen/gexch/common"
	"errors"
	"strings"
)

const (
	OKExWebSocketBaseApi   = "wss://real.okex.com:10440/websocket/okexapi"
	OKExWebSocketTicker    = "ok_sub_futureusd_%s_ticker_%s"
	OKExWebSocketDepthIncr = "ok_sub_future_%s_depth_%s_usd"
	OKExWebSocketDepthFull = "ok_sub_futureusd_%s_depth_%s_%d"
	OKExWebSocketTrade     = "ok_sub_futureusd_%s_trade_%s"
)

func (o *OKEx) WsConnect() (error) {
	var Dialer websocket.Dialer
	var err error
	o.WebsocketConn, _, err = Dialer.Dial(OKExWebSocketBaseApi, http.Header{})

	if err != nil {
		log.Printf("%s Unable to connect to Websocket. Error: %s\n", o.Name, err)
		return err
	}

	if o.Verbose {
		log.Printf("%s Connected to Websocket.\n", o.Name)
	}

	o.WebsocketConn.SetPingHandler(o.PingHandler)
	o.WebsocketConn.SetCloseHandler(o.CloseHandler)
	return nil
}

func (o *OKEx) CloseHandler(code int, text string) error {
	log.Println("closed")
	return nil
}

func (o *OKEx) PingHandler(message string) error {
	err := o.WebsocketConn.WriteControl(websocket.PingMessage, []byte("{'event':'ping'}"), time.Now().Add(time.Second))

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (o *OKEx) AddChannel(channel string) {
	event := OKCoinWebsocketEvent{"addChannel", channel}
	json, err := common.JSONEncode(event)
	if err != nil {
		log.Println(err)
		return
	}
	err = o.WebsocketConn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		log.Println(err)
		return
	}

	if o.Verbose {
		log.Printf("%s Adding channel: %s\n", o.Name, channel)
	}
}

func (o *OKEx) RemoveChannel(channel string) {
	event := OKCoinWebsocketEvent{"removeChannel", channel}
	json, err := common.JSONEncode(event)
	if err != nil {
		log.Println(err)
		return
	}
	err = o.WebsocketConn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		log.Println(err)
		return
	}

	if o.Verbose {
		log.Printf("%s Removing channel: %s\n", o.Name, channel)
	}
}

func (o *OKEx) WsProcessResponse(resp interface{}) (interface{}, error) {
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
			log.Printf("%s Websocket: channel %s error: %s.", o.Name, channel, WebsocketErrors[errorCodeStr])
		}
		return nil, errors.New(fmt.Sprintf("channel %s error: %s", channel, WebsocketErrors[errorCodeStr]))
	}

	if success == "true" && data == nil {
		return nil, errors.New("no data found")
	}

	//fmt.Println(channel)
	//fmt.Println(string(data))

	switch true {
	case channel == "addChannel":
		return channel, nil
	case strings.Contains(channel, "ok_sub_futureusd"):
		if strings.Contains(channel, "ticker") {
			v := OKCoinWebsocketFuturesTicker{}
			err := common.JSONDecode(data, &v)
			return v, err
		} else if strings.Contains(channel, "depth") {
			v := OKCoinWebsocketOrderbook{}
			err := common.JSONDecode(data, &v)
			return v, err
		} else if strings.Contains(channel, "trade") {
			v := [][]string{}
			err := common.JSONDecode(data, &v)
			return v, err
		} else {
			goto NoChannelTypeMatched
		}
	case strings.Contains(channel, "ok_sub_future") && strings.Contains(channel, "depth"):
		v := OKCoinWebsocketOrderbook{}
		err := common.JSONDecode(data, &v)
		return v, err
	default:
		goto NoChannelTypeMatched
	}

NoChannelTypeMatched:
	return channel, errors.New(fmt.Sprintf("unknown channel type: %s", channel))
}

var (
	WebsocketErrors = map[string]string{
		"10001": "Illegal parameters",
		"10002": "Authentication failure",
		"10003": "This connection has requested other user data",
		"10004": "This connection did not request this user data",
		"10005": "System error",
		"10009": "Order does not exist",
		"10010": "Insufficient funds",
		"10011": "Order quantity too low",
		"10012": "Only support btc_usd/btc_cny ltc_usd/ltc_cny",
		"10014": "Order price must be between 0 - 1,000,000",
		"10015": "Channel subscription temporally not available",
		"10016": "Insufficient coins",
		"10017": "WebSocket authorization error",
		"10100": "User frozen",
		"10216": "Non-public API",
		"20001": "User does not exist",
		"20002": "User frozen",
		"20003": "Frozen due to force liquidation",
		"20004": "Future account frozen",
		"20005": "User future account does not exist",
		"20006": "Required field can not be null",
		"20007": "Illegal parameter",
		"20008": "Future account fund balance is zero",
		"20009": "Future contract status error",
		"20010": "Risk rate information does not exist",
		"20011": `Risk rate bigger than 90% before opening position`,
		"20012": `Risk rate bigger than 90% after opening position`,
		"20013": "Temporally no counter party price",
		"20014": "System error",
		"20015": "Order does not exist",
		"20016": "Liquidation quantity bigger than holding",
		"20017": "Not authorized/illegal order ID",
		"20018": `Order price higher than 105% or lower than 95% of the price of last minute`,
		"20019": "IP restrained to access the resource",
		"20020": "Secret key does not exist",
		"20021": "Index information does not exist",
		"20022": "Wrong API interface",
		"20023": "Fixed margin user",
		"20024": "Signature does not match",
		"20025": "Leverage rate error",
	}
)
