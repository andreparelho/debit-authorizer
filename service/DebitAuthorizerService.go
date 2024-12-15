package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	request "github.com/andreparelho/debit-authorizer/model/common"
	serviceDTO "github.com/andreparelho/debit-authorizer/model/service"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

const LAST_FIVE_MINUTES = 5 * time.Minute
const MAX_TOTAL_AMOUNT = 1000
const EMPTY_VALUE = ""

var transactionHitorical = make(map[string]serviceDTO.Client)
var mutex sync.Mutex
var message []byte
var valueMessage string

func DebitAuthorizerService(request request.RequestAuthorizerDebit) ([]byte, error) {
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
	if !isCreated {
		transactionHitorical[request.ClientId] = serviceDTO.Client{
			LastPayment: dateTime,
			TotalAmount: request.Amount,
		}
		logger.ServiceLoggerInfo(client, clientId, "client created")
	}

	var totalAmount = client.TotalAmount + request.Amount
	if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		valueMessage = "sorry you have reached your debit limit"
		message = []byte(fmt.Sprintf(`{"message": "%s"}`, valueMessage))
		var errorMessage error = errors.New(valueMessage)

		logger.ServiceLoggerError(client, clientId, valueMessage)
		return message, errorMessage
	}

	if totalAmount > MAX_TOTAL_AMOUNT {
		valueMessage = "sorry the amount sent is greater than the allowed limit"
		message = []byte(fmt.Sprintf(`{"message": "%s"}`, valueMessage))
		var errorMessage error = errors.New(valueMessage)

		logger.ServiceLoggerError(client, clientId, valueMessage)
		return message, errorMessage
	}

	client.LastPayment = dateTime
	client.TotalAmount = totalAmount
	transactionHitorical[clientId] = client

	message = []byte(`{"message": "debit approved"}`)
	logger.ServiceLoggerInfo(client, clientId, "debit approved")
	return message, nil
}
