package service

import (
	"errors"
	"sync"
	"time"

	request "github.com/andreparelho/debit-authorizer/model/common"
	serviceDTO "github.com/andreparelho/debit-authorizer/model/service"
	logger "github.com/andreparelho/debit-authorizer/util/logUtil"
)

const LAST_FIVE_MINUTES = 5 * time.Minute
const MAX_TOTAL_AMOUNT = 1000
const EMPTY_VALUE = ""

var clientHistorical = make(map[string]serviceDTO.Client)
var mutex sync.Mutex
var message []byte

func DebitAuthorizerService(request request.RequestAuthorizerDebit) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var now time.Time = time.Now()
	var dateTime time.Time
	if request.DateTime.IsZero() {
		dateTime = now
	}

	var clientId = request.ClientId
	client, isCreated := clientHistorical[clientId]
	if !isCreated {
		clientHistorical[request.ClientId] = serviceDTO.Client{
			LastPayment: dateTime,
			TotalAmount: request.Amount,
		}
		logger.ServiceLoggerInfo(client, clientId, "client created")
	}

	var totalAmount = client.TotalAmount + request.Amount
	if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		message = []byte(`{"message": "Sorry you have reached your debit limit"}`)
		var errorMessage error = errors.New("Sorry you have reached your debit limit")

		logger.ServiceLoggerInfo(client, clientId, "Sorry you have reached your debit limit")
		return message, errorMessage
	}

	client.LastPayment = dateTime
	client.TotalAmount = totalAmount
	clientHistorical[clientId] = client

	message = []byte(`{"message": "debit approved"}`)
	logger.ServiceLoggerInfo(client, clientId, "debit approved")
	return message, nil
}
