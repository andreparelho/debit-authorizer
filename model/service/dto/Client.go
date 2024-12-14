package model

import "time"

type Client struct {
	LastPayment time.Time
	TotalAmount float64
}
