package repository

import (
	"time"
)

type ClientHistorical interface {
	CreateClientHistorical(cl Client, dateTime time.Time, amount float64)
	UpdateClientHistorical(cl Client, id string, dateTime time.Time, totalAmount float64, amountRequest float64)
	GetClientHistorical(id string) Client
}

func NewClientHistorical(transactions map[string]Client) *client {
	return &client{
		Transactions: transactions,
	}
}

func (c *client) CreateClientHistorical(cl Client, dateTime time.Time, amount float64) {
	cl.Historical = append(cl.Historical, Historical{
		Amount:   amount,
		DateTime: dateTime,
	})
	c.Transactions[cl.ClientId] = cl
}

func (c *client) UpdateClientHistorical(cl Client, id string, dateTime time.Time, totalAmount float64, amountRequest float64) {
	cl.LastPayment = dateTime
	cl.TotalAmount = totalAmount
	cl.Historical = append(cl.Historical, Historical{
		Amount:   amountRequest,
		DateTime: dateTime,
	})

	c.Transactions[id] = cl
}

func (c *client) GetClientHistorical(id string) Client {
	return c.Transactions[id]
}
