package model

type RequestAuthorizerDebitHandler struct {
	ClientId string  `json:"clientId"`
	Amount   float64 `json:"amount"`
}
