package handler

import (
	"time"

	"github.com/andreparelho/debit-authorizer/internal/repository"
)

type RequestAuthorizerDebit struct {
	ClientId string    `json:"clientId"`
	DateTime time.Time `json:"dateTime"`
	Amount   float64   `json:"amount"`
}

type ResponseAuthorizerDebit struct {
	Message          string            `json:"message"`
	ClientHistorical repository.Client `json:"transactions"`
}
