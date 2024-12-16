package handler

import (
	"encoding/json"
	"net/http"

	requestHandler "github.com/andreparelho/debit-authorizer/model/common"
	service "github.com/andreparelho/debit-authorizer/service"
	httpUtil "github.com/andreparelho/debit-authorizer/util/httpUtil"
)

const POST_METHOD string = "POST"
const EMPTY_VALUE string = ""
const MINIMUM_AMOUNT_REQUEST float64 = 0.01

func DebitAuthorizerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var httpUtil *httpUtil.HttpUtil = httpUtil.ResponseJSONConstructor(responseWriter)

	if request.Method != POST_METHOD {
		var message []byte = []byte(`{"message": "invalid method"}`)
		httpUtil.ResponseJSON(message, http.StatusMethodNotAllowed)
		return
	}

	var requestHandler requestHandler.RequestAuthorizerDebit
	var requestHandlerError error = json.NewDecoder(request.Body).Decode(&requestHandler)
	if requestHandlerError != nil {
		var message []byte = []byte(`{"message": "error to decode json"}`)
		httpUtil.ResponseJSON(message, http.StatusInternalServerError)
		return
	}

	var clientId string = requestHandler.ClientId
	if clientId == EMPTY_VALUE {
		var message []byte = []byte(`{"message": "propertie clientId is empty"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	var amountRequest float64 = requestHandler.Amount
	if amountRequest <= MINIMUM_AMOUNT_REQUEST {
		var message []byte = []byte(`{"message": "propertie amount is empty or less than minimum"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	response, errorService := service.DebitAuthorizerService(requestHandler)
	if errorService != nil {
		var message []byte = []byte(`{"message": "sorry you have reached your debit limit"}`)
		httpUtil.ResponseJSON(message, http.StatusTooManyRequests)
		return
	}

	httpUtil.ResponseJSONSuccess(response, http.StatusOK)
}
