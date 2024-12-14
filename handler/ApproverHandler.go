package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	requestHandler "github.com/andreparelho/debit-authorizer/model/handler/request"
	service "github.com/andreparelho/debit-authorizer/service"
	httpUtil "github.com/andreparelho/debit-authorizer/util/httpUtil"
)

const POST_METHOD string = "POST"
const EMPTY string = ""
const ZERO int64 = 0

func ApproverHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var httpUtil = httpUtil.ResponseJSONConstructor(responseWriter)

	if request.Method != POST_METHOD {
		var message []byte = []byte(`{"message": "Invalid method"}`)
		httpUtil.ResponseJSON(message, http.StatusMethodNotAllowed)
		responseWriter.Write(message)
		return
	}

	var requestHandler requestHandler.RequestApproverHandler
	var requestHandlerError error = json.NewDecoder(request.Body).Decode(&requestHandler)
	if requestHandlerError != nil {
		var message []byte = []byte(`{"message": "Error to decode json"}`)
		httpUtil.ResponseJSON(message, http.StatusInternalServerError)
		return
	}

	var clientId string = requestHandler.ClientId
	if clientId == EMPTY {
		var message []byte = []byte(`{"message": "Propertie clientId is empty"}`)
		responseWriter.Header().Set("Content-Type", "application/json")
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	var amountRequest int64 = requestHandler.Amount
	if amount := strconv.Itoa(int(amountRequest)); amountRequest <= ZERO || amount == EMPTY {
		var message []byte = []byte(`{"message": "Propertie amout is empty or less than zero"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	defer request.Body.Close()

	message, errorService := service.ApproverService(requestHandler)
	if errorService != nil {
		httpUtil.ResponseJSON(message, http.StatusTooManyRequests)
		return
	}

	httpUtil.ResponseJSON(message, http.StatusOK)
}
