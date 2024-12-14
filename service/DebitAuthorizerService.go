package service

import (
	"errors"
	"sync"
	"time"

	requestHandler "github.com/andreparelho/debit-authorizer/model/handler/request"
	serviceDTO "github.com/andreparelho/debit-authorizer/model/service/dto"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

const LAST_FIVE_MINUTES = 1 * time.Minute

var clientHistory = make(map[string]serviceDTO.Client)
var mutex sync.Mutex

func DebitAuthorizerService(requestHandler requestHandler.RequestApproverHandler) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var now time.Time = time.Now()
	var clientId = requestHandler.ClientId

	client, isCreated := clientHistory[clientId]
	if !isCreated {
		clientHistory[requestHandler.ClientId] = serviceDTO.Client{
			LastPayment: now,
			TotalAmount: requestHandler.Amount,
		}
		logger.ServiceLoggerInfo(client, clientId, "client created")
	}

	var totalAmount = client.TotalAmount + requestHandler.Amount
	if totalAmount > 1000 && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		var message []byte = []byte(`{"message": "Sorry you have reached your debit limit"}`)
		logger.ServiceLoggerInfo(client, clientId, "Sorry you have reached your debit limit")

		var errorMessage error = errors.New("Sorry you have reached your debit limit")

		return message, errorMessage
	}

	client.LastPayment = now
	client.TotalAmount = totalAmount
	clientHistory[clientId] = client

	var message []byte = []byte(`{"message": "debit approved"}`)
	logger.ServiceLoggerInfo(client, clientId, "debit approved")
	return message, nil
}
