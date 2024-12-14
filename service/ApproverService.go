package service

import (
	"sync"
	"time"

	requestHandler "github.com/andreparelho/debit-authorizer/model/handler/request"
	serviceDTO "github.com/andreparelho/debit-authorizer/model/service/dto"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

const LAST_FIVE_MINUTES = 5 * time.Minute

var clientHistory = make(map[string]serviceDTO.Client)
var mutex sync.Mutex
var now time.Time = time.Now()

func ApproverService(requestHandler requestHandler.RequestApproverHandler) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	client, isCreated := clientHistory[requestHandler.ClientId]
	if !isCreated {
		clientHistory[requestHandler.ClientId] = serviceDTO.Client{
			LastPayment: now,
			TotalAmount: requestHandler.Amount,
		}
		logger.ServiceLoggerInfo(client, requestHandler.ClientId, "client created")
	}

	var totalAmount = client.TotalAmount + requestHandler.Amount
	if totalAmount > 1000 && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		var message []byte = []byte(`{"message": "Sorry you have reached your debit limit"}`)
		logger.ServiceLoggerInfo(client, requestHandler.ClientId, "Sorry you have reached your debit limit")
		return message, nil
	}

	client.LastPayment = now
	client.TotalAmount = totalAmount
	clientHistory[requestHandler.ClientId] = client

	var message []byte = []byte(`{"message": "debit approved"}`)
	logger.ServiceLoggerInfo(client, requestHandler.ClientId, "debit approved")
	return message, nil
}
