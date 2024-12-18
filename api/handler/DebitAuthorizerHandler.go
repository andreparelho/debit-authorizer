package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/andreparelho/debit-authorizer/model/common"
	service "github.com/andreparelho/debit-authorizer/service"
	httpUtil "github.com/andreparelho/debit-authorizer/util/httpUtil"
)

const POST_METHOD string = "POST"
const EMPTY_VALUE string = ""
const MINIMUM_AMOUNT_REQUEST float64 = 0.01

var message []byte

func DebitAuthorizerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var httpUtil *httpUtil.HttpUtil = httpUtil.ResponseJSONConstructor(responseWriter)

	if request.Method != POST_METHOD {
		message = []byte(`{"message": "invalid method"}`)
		httpUtil.ResponseJSON(message, http.StatusMethodNotAllowed)
		return
	}

	var requestHandler common.RequestAuthorizerDebit
	var requestHandlerError error = json.NewDecoder(request.Body).Decode(&requestHandler)
	if requestHandlerError != nil {
		message = []byte(`{"message": "error to decode json"}`)
		httpUtil.ResponseJSON(message, http.StatusInternalServerError)
		return
	}

	var clientId string = requestHandler.ClientId
	if clientId == EMPTY_VALUE {
		message = []byte(`{"message": "propertie clientId is empty"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	var amountRequest float64 = requestHandler.Amount
	if amountRequest <= MINIMUM_AMOUNT_REQUEST {
		message = []byte(`{"message": "propertie amount is empty or less than minimum"}`)
		httpUtil.ResponseJSON(message, http.StatusBadRequest)
		return
	}

	responseService, errorService := service.DebitAuthorizerService(requestHandler)
	if errorService != nil {
		var errorMessage string = fmt.Sprintf(`{"message": "%s"}`, errorService.Error())
		message = []byte(errorMessage)
		httpUtil.ResponseJSON(message, http.StatusTooManyRequests)
		return
	}

	var response common.ResponseAuthorizerDebit = common.ResponseAuthorizerDebit{
		Message:          "debit authorized",
		ClientHistorical: responseService,
	}

	httpUtil.ResponseJSONSuccess(response, http.StatusOK)
}
