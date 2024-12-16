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

var transactionHitorical = make(map[string]model.Client)
var mutex sync.Mutex
var valueMessage string

func DebitAuthorizerService(request request.RequestAuthorizerDebit) (model.Client, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var now time.Time = time.Now()
	var dateTime time.Time
	if request.DateTime.IsZero() {
		dateTime = now
	} else {
		dateTime = request.DateTime
	}

	var clientId = request.ClientId
	client, isCreated := transactionHitorical[clientId]
	var totalAmount = client.TotalAmount + request.Amount

	if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		valueMessage = "sorry you have reached your debit limit"
		var errorMessage error = errors.New(valueMessage)

		logger.ServiceLoggerError(client.LastPayment, client.TotalAmount, clientId, valueMessage)
		return model.Client{}, errorMessage
	}

	if totalAmount > MAX_TOTAL_AMOUNT {
		valueMessage = "sorry the amount sent is greater than the allowed limit"
		var errorMessage error = errors.New(valueMessage)

		logger.ServiceLoggerError(client.LastPayment, client.TotalAmount, client.ClientId, valueMessage)
		return model.Client{}, errorMessage
	}

	if !isCreated {
		client = model.Client{
			ClientId:    clientId,
			LastPayment: dateTime,
			TotalAmount: request.Amount,
			Historical:  []model.Historical{},
		}

		repository.CreateClientHistorical(transactionHitorical, client, dateTime, request.Amount)
		logger.ServiceLoggerInfo(dateTime, request.Amount, clientId, "client created")
	}

	if isCreated {
		repository.UpdateClientHistorical(client, transactionHitorical, clientId, dateTime, totalAmount, request.Amount)
	}

	valueMessage = "debit approved"
	logger.ServiceLoggerInfo(client.LastPayment, client.TotalAmount, clientId, valueMessage)

	return repository.GetClientHitorical(clientId, transactionHitorical), nil
}
