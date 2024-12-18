package model

import "time"

type Client struct {
	ClientId    string       `json:"clientId"`
	LastPayment time.Time    `json:"lastPayment"`
	TotalAmount float64      `json:"totalAmount"`
	Historical  []Historical `json:"historical"`
}

type Historical struct {
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"dateTime"`
}
