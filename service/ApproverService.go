package service

import (
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	requestHandler "github.com/andreparelho/credit-approver/model/handler/request"
	serviceDTO "github.com/andreparelho/credit-approver/model/service/dto"
)

const LAST_FIVE_MINUTES = 5 * time.Minute

var clientHistory = make(map[string]serviceDTO.Client)
var mutex sync.Mutex
var now time.Time = time.Now()
var jsonHandler = slog.NewJSONHandler(os.Stderr, nil)
var myslog = slog.New(jsonHandler)

type Service struct {
	ResponseWriter http.ResponseWriter
}

func ServiceApproverImpl(responseWriter http.ResponseWriter) *Service {
	return &Service{
		ResponseWriter: responseWriter,
	}
}

func (writer *Service) ApproverService(requestHandler requestHandler.RequestApproverHandler) {
	mutex.Lock()
	defer mutex.Unlock()

	client, isCreated := clientHistory[requestHandler.ClientId]
	if !isCreated {
		clientHistory[requestHandler.ClientId] = serviceDTO.Client{
			LastPayment: now,
			TotalAmount: requestHandler.Amount,
		}
	}

	if now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
		var message []byte = []byte(`{"message": "Please wait 5 minutes"}`)
		writer.ResponseWriter.Header().Set("Content-Type", "application/json")
		writer.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
		writer.ResponseWriter.Write(message)
		return
	}

	var totalAmount = client.TotalAmount + requestHandler.Amount
	if totalAmount > 1000 {
		var message []byte = []byte(`{"message": "Sorry you have reached your credit limit"}`)
		writer.ResponseWriter.Header().Set("Content-Type", "application/json")
		writer.ResponseWriter.WriteHeader(http.StatusBadRequest)
		writer.ResponseWriter.Write(message)
		return
	}

	client.LastPayment = now
	client.TotalAmount = totalAmount
	clientHistory[requestHandler.ClientId] = client

	var message []byte = []byte(`{"message": "Credit approved"}`)
	writer.ResponseWriter.Header().Set("Content-Type", "application/json")
	writer.ResponseWriter.WriteHeader(http.StatusOK)
	writer.ResponseWriter.Write(message)
}
