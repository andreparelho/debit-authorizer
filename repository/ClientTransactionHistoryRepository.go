package repository

import (
	"time"

	model "github.com/andreparelho/debit-authorizer/model/service"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

func CreateClientHistorical(clientHistorical map[string]model.Client, client model.Client, dateTime time.Time, amount float64) {

	var transacation model.Historical = model.Historical{
		Amount:   amount,
		DateTime: dateTime,
	}

	client.Historical = append(client.Historical, transacation)
	clientHistorical[client.ClientId] = client

	historical := clientHistorical[client.ClientId]
	logger.RepositoryLoggerInfo(client.ClientId, historical.Historical, "created client")
}

func UpdateClientHistorical(client model.Client, clientHistorical map[string]model.Client, clientId string, dateTime time.Time, totalAmount float64, amountRequest float64) {
	var transacation model.Historical = model.Historical{
		Amount:   amountRequest,
		DateTime: dateTime,
	}

	client.LastPayment = dateTime
	client.TotalAmount = totalAmount
	client.Historical = append(client.Historical, transacation)
	clientHistorical[clientId] = client

	historical := clientHistorical[clientId]

	logger.RepositoryLoggerInfo(clientId, historical.Historical, "updated client")
}

func GetClientHitorical(clientId string, clientHistorical map[string]model.Client) model.Client {
	return clientHistorical[clientId]
}
