package common

import "net/http"

type ExchangeBase struct {
	Name       string
	Enable     bool
	BaseUri    string
	HttpClient *http.Client
	ApiKey     string
	SecretKey  string
	TakerFee   float64
	MakerFee   float64
}
