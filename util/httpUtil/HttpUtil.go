package util

import "net/http"

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
