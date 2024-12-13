package model

type RequestApproverHandler struct {
	ClientId string `json:"clientId"`
	Amount   int64  `json:"amount"`
}
