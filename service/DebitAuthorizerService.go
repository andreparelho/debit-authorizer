package service

import (
	"errors"
	"sync"
	"time"

	request "github.com/andreparelho/debit-authorizer/model/common"
	model "github.com/andreparelho/debit-authorizer/model/service"
	repository "github.com/andreparelho/debit-authorizer/repository"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

const LAST_FIVE_MINUTES = 5 * time.Minute
const MAX_TOTAL_AMOUNT = 1000
const EMPTY_VALUE = ""

var transactionHistorical = make(map[string]model.Client)
var mutex sync.Mutex
var valueMessage string

func DebitAuthorizerService(request request.RequestAuthorizerDebit) (model.Client, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var now time.Time = time.Now()
	var dateTime time.Time = getDate(request.DateTime, now)

	var clientId = request.ClientId
	client, isCreated := transactionHistorical[clientId]

	var totalAmount = client.TotalAmount + request.Amount
	if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		valueMessage = "sorry you have reached your debit limit"
		var errorMessage error = errors.New(valueMessage)

		logger.ServiceLoggerError(clientId, request.Amount, totalAmount, errorMessage.Error())
		return model.Client{}, errorMessage
	}

	if totalAmount > MAX_TOTAL_AMOUNT {
		valueMessage = "sorry the amount sent is greater than the allowed limit"
		var errorMessage error = errors.New(valueMessage)

		logger.ServiceLoggerError(clientId, request.Amount, totalAmount, errorMessage.Error())
		return model.Client{}, errorMessage
	}

	validateClient(isCreated, client, clientId, dateTime, request.Amount, totalAmount)

	valueMessage = "debit authorized"
	logger.ServiceLoggerInfo(clientId, client.LastPayment, totalAmount, valueMessage)

	var response model.Client = repository.GetClientHitorical(clientId, transactionHistorical)
	return response, nil
}

func getDate(requestDate time.Time, now time.Time) time.Time {
	if requestDate.IsZero() {
		return now
	}

	return requestDate
}

func validateClient(isCreated bool, client model.Client, clientId string, dateTime time.Time, amount float64, totalAmount float64) {
	if !isCreated {
		client = model.Client{
			ClientId:    clientId,
			LastPayment: dateTime,
			TotalAmount: amount,
			Historical:  []model.Historical{},
		}
		repository.CreateClientHistorical(transactionHistorical, client, dateTime, amount)
		logger.ServiceLoggerInfo(clientId, client.LastPayment, totalAmount, "client created")
	} else {
		repository.UpdateClientHistorical(client, transactionHistorical, clientId, dateTime, totalAmount, amount)
		logger.ServiceLoggerInfo(clientId, client.LastPayment, totalAmount, "client updated")
	}
}
