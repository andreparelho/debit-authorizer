package model

import model "github.com/andreparelho/debit-authorizer/model/service"

type ResponseAuthorizerDebit struct {
	Message          string       `json:"message"`
	ClientHistorical model.Client `json:"transactions"`
}
