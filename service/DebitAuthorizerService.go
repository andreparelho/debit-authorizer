package service

import (
	"errors"
	"sync"
	"time"

	requestHandler "github.com/andreparelho/debit-authorizer/model/handler/request"
	serviceDTO "github.com/andreparelho/debit-authorizer/model/service/dto"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

const LAST_FIVE_MINUTES = 5 * time.Minute
const MAX_TOTAL_AMOUNT = 1000

var clientHistorical = make(map[string]serviceDTO.Client)
var mutex sync.Mutex
var message []byte

func DebitAuthorizerService(requestHandler requestHandler.RequestAuthorizerDebitHandler) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var now time.Time = time.Now()
	var clientId = requestHandler.ClientId

	client, isCreated := clientHistorical[clientId]
	if !isCreated {
		clientHistorical[requestHandler.ClientId] = serviceDTO.Client{
			LastPayment: now,
			TotalAmount: requestHandler.Amount,
		}
		logger.ServiceLoggerInfo(client, clientId, "client created")
	}

	var totalAmount = client.TotalAmount + requestHandler.Amount
	if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		message = []byte(`{"message": "Sorry you have reached your debit limit"}`)
		logger.ServiceLoggerInfo(client, clientId, "Sorry you have reached your debit limit")

		var errorMessage error = errors.New("Sorry you have reached your debit limit")

		return message, errorMessage
	}

	client.LastPayment = now
	client.TotalAmount = totalAmount
	clientHistorical[clientId] = client

	message = []byte(`{"message": "debit approved"}`)
	logger.ServiceLoggerInfo(client, clientId, "debit approved")
	return message, nil
}
