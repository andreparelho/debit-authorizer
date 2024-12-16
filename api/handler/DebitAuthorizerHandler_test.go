package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreparelho/debit-authorizer/api/handler"
	"github.com/stretchr/testify/assert"
)

const ENDPOINT_URL string = "http://localhost:8080/authorizer-debit"

func TestDebitAuthorizerHandler(test *testing.T) {
	test.Run("Deve validar erro quando o metodo http diferente de POST", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": "1","amount": 500}`))
		var request *http.Request = httptest.NewRequest("GET", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var errorMessage string = `{"message": "invalid method"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusMethodNotAllowed)
		assert.Equal(test, responseWriter.Body.String(), errorMessage)
	})

	test.Run("Deve validar erro quando o request invalido", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": 1,"amount": "500"}`))
		var request *http.Request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var errorMessage string = `{"message": "error to decode json"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusInternalServerError)
		assert.Equal(test, responseWriter.Body.String(), errorMessage)
	})

	test.Run("Deve validar erro quando a proriedade clientId estiver vazia", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": "","amount": 500}`))
		var request *http.Request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var errorMessage string = `{"message": "propertie clientId is empty"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusBadRequest)
		assert.Equal(test, responseWriter.Body.String(), errorMessage)
	})

	test.Run("Deve validar erro quando a proriedade amout abaixo do valor minimo", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": "1","amount": 0.00}`))
		var request *http.Request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var errorMessage string = `{"message": "propertie amount is empty or less than minimum"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusBadRequest)
		assert.Equal(test, responseWriter.Body.String(), errorMessage)
	})

	test.Run("Deve validar erro quando enviar valor acima do esperado no primeiro request", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": "1","amount": 10000}`))
		var request *http.Request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var errorMessage string = `{"message": "sorry the amount sent is greater than the allowed limit"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusTooManyRequests)
		assert.Equal(test, responseWriter.Body.String(), errorMessage)
	})

	test.Run("Deve validar erro quando atingir o limite", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": "1","amount": 1000}`))
		var request *http.Request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		responseWriter = httptest.NewRecorder()
		bodyRequest = bytes.NewBuffer([]byte(`{"clientId": "1","amount": 1}`))
		request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var errorMessage string = `{"message": "sorry you have reached your debit limit"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusTooManyRequests)
		assert.Equal(test, responseWriter.Body.String(), errorMessage)
	})

	test.Run("Deve validar sucesso quando enviar os parametros corretos", func(test *testing.T) {
		var responseWriter *httptest.ResponseRecorder = httptest.NewRecorder()
		var bodyRequest *bytes.Buffer = bytes.NewBuffer([]byte(`{"clientId": "2","amount": 1000}`))
		var request *http.Request = httptest.NewRequest("POST", ENDPOINT_URL, bodyRequest)

		handler.DebitAuthorizerHandler(responseWriter, request)

		var successMessage string = `{"message": "debit approved"}`

		assert.True(test, responseWriter.Result().StatusCode == http.StatusOK)
		assert.Equal(test, responseWriter.Body.String(), successMessage)
	})
}
