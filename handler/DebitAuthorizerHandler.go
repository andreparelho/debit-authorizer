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
const MINIMUM_AMOUNT_REQUEST float64 = 0.01

func DebitAuthorizerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var httpUtil *httpUtil.HttpUtil = httpUtil.ResponseJSONConstructor(responseWriter)

	if request.Method != POST_METHOD {
		var message []byte = []byte(`{"message": "Invalid method"}`)
		httpUtil.ResponseJSON(message, http.StatusMethodNotAllowed)
		return
	}

	var requestHandler requestHandler.RequestAuthorizerDebitHandler
	var requestHandlerError error = json.NewDecoder(request.Body).Decode(&requestHandler)
	if requestHandlerError != nil {
		var message []byte = []byte(`{"message": "Error to decode json"}`)
		httpUtil.ResponseJSON(message, http.StatusInternalServerError)
		return
	}

	var clientId string = requestHandler.ClientId
	if clientId == EMPTY {
		var message []byte = []byte(`{"message": "Propertie clientId is empty"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	var amountRequest float64 = requestHandler.Amount
	if amountString := strconv.Itoa(int(amountRequest)); amountString == EMPTY || amountRequest <= MINIMUM_AMOUNT_REQUEST {
		var message []byte = []byte(`{"message": "Propertie amount is empty or less than minimum"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	message, errorService := service.DebitAuthorizerService(requestHandler)
	if errorService != nil {
		httpUtil.ResponseJSON(message, http.StatusTooManyRequests)
		return
	}

	httpUtil.ResponseJSON(message, http.StatusOK)
}
