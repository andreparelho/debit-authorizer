package util

import (
	"encoding/json"
	"net/http"

	response "github.com/andreparelho/debit-authorizer/model/common"
)

type HttpUtil struct {
	ResponseWriter http.ResponseWriter
}

func ResponseJSONConstructor(responseWriter http.ResponseWriter) *HttpUtil {
	return &HttpUtil{
		ResponseWriter: responseWriter,
	}
}

func (writer *HttpUtil) ResponseJSON(message []byte, statusCode int) {
	writer.ResponseWriter.Header().Set("Content-Type", "application/json")
	writer.ResponseWriter.WriteHeader(statusCode)
	writer.ResponseWriter.Write(message)
}

func (writer *HttpUtil) ResponseJSONSuccess(responseService response.ResponseAuthorizerDebit, statusCode int) {
	response, _ := json.Marshal(responseService)
	writer.ResponseWriter.Header().Set("Content-Type", "application/json")
	writer.ResponseWriter.WriteHeader(statusCode)
	writer.ResponseWriter.Write(response)
}
