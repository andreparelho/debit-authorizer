package model

type RequestApproverHandler struct {
	ClientId string  `json:"clientId"`
	Amount   float64 `json:"amount"`
}
