package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	requestHandler "github.com/andreparelho/credit-approver/model/handler/request"
	service "github.com/andreparelho/credit-approver/service"
)

const POST_METHOD string = "POST"
const EMPTY string = ""
const ZERO int64 = 0

func ApproverHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != POST_METHOD {
		var message []byte = []byte(`{"message": "Invalid method"}`)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		responseWriter.Write(message)
		return
	}

	var requestHandler requestHandler.RequestApproverHandler
	var requestHandlerError error = json.NewDecoder(request.Body).Decode(&requestHandler)
	if requestHandlerError != nil {
		var message []byte = []byte(`{"message": "Error to decode json"}`)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write(message)
		return
	}

	var clientId string = requestHandler.ClientId
	if clientId == EMPTY {
		var message []byte = []byte(`{"message": "Propertie clientId is empty"}`)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(message)
		return
	}

	var amountRequest int64 = requestHandler.Amount
	if amount := strconv.Itoa(int(amountRequest)); amountRequest <= ZERO || amount == EMPTY {
		var message []byte = []byte(`{"message": "Propertie amout is empty or less than zero"}`)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(message)
		return
	}

	defer request.Body.Close()

	service := service.ServiceApproverImpl(responseWriter)
	service.ApproverService(requestHandler)
}
