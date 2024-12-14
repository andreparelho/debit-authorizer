package model

import "time"

type RequestAuthorizerDebit struct {
	ClientId string    `json:"clientId"`
	DateTime time.Time `json:"dateTime"`
	Amount   float64   `json:"amount"`
}
