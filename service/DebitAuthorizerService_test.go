package service_test

import (
	"testing"
	"time"

	request "github.com/andreparelho/debit-authorizer/model/common"
	client "github.com/andreparelho/debit-authorizer/model/service"
	"github.com/andreparelho/debit-authorizer/service"
	"github.com/stretchr/testify/assert"
)

func TestDebitAuthorizerService(test *testing.T) {

	test.Run("Deve validar com sucesso o response", func(test *testing.T) {
		var dateTime = time.Now()

		var request = request.RequestAuthorizerDebit{
			ClientId: "1",
			DateTime: dateTime,
			Amount:   500,
		}

		response, errorService := service.DebitAuthorizerService(request)

		var historicalClientCreated = client.Historical{
			Amount:   500,
			DateTime: dateTime,
		}

		var clientCreated client.Client = client.Client{
			ClientId:    "1",
			LastPayment: dateTime,
			TotalAmount: 500,
			Historical:  append(response.Historical, historicalClientCreated),
		}

		assert.Nil(test, errorService)
		assert.Equal(test, clientCreated.ClientId, response.ClientId)
		assert.Equal(test, clientCreated.LastPayment, response.LastPayment)
		assert.Equal(test, clientCreated.TotalAmount, response.TotalAmount)

		for _, responseValue := range response.Historical {
			assert.Equal(test, responseValue, historicalClientCreated)
		}
	})

	test.Run("Deve validar erro quando enviar o segundo request por ser acima do limite permitido", func(test *testing.T) {
		var dateTime = time.Now()

		var firstRequest = request.RequestAuthorizerDebit{
			ClientId: "2",
			DateTime: dateTime,
			Amount:   1000,
		}

		responseFirstRequest, errorFirstRequest := service.DebitAuthorizerService(firstRequest)

		var historicalClientCreated = client.Historical{
			Amount:   1000,
			DateTime: dateTime,
		}

		var clientCreated client.Client = client.Client{
			ClientId:    "2",
			LastPayment: dateTime,
			TotalAmount: 1000,
			Historical:  append(responseFirstRequest.Historical, historicalClientCreated),
		}

		assert.Nil(test, errorFirstRequest)
		assert.Equal(test, clientCreated.ClientId, responseFirstRequest.ClientId)
		assert.Equal(test, clientCreated.LastPayment, responseFirstRequest.LastPayment)
		assert.Equal(test, clientCreated.TotalAmount, responseFirstRequest.TotalAmount)

		for _, responseValue := range responseFirstRequest.Historical {
			assert.Equal(test, responseValue, historicalClientCreated)
		}

		var secondRequest = request.RequestAuthorizerDebit{
			ClientId: "2",
			DateTime: time.Now(),
			Amount:   100,
		}

		responseSecondRequest, errorSecondRequest := service.DebitAuthorizerService(secondRequest)
		var messageExpected string = "sorry you have reached your debit limit"

		assert.NotNil(test, errorSecondRequest)
		assert.EqualErrorf(test, errorSecondRequest, errorSecondRequest.Error(), "sorry you have reached your debit limit")
		assert.Equal(test, messageExpected, errorSecondRequest.Error())
		assert.Empty(test, responseSecondRequest)
	})

	test.Run("Deve validar erro quando o valor enviado estiver acima do limite maxima por transacao", func(test *testing.T) {
		var request = request.RequestAuthorizerDebit{
			ClientId: "3",
			DateTime: time.Now(),
			Amount:   1500,
		}

		responseRequest, errorRequest := service.DebitAuthorizerService(request)
		var messageExpected string = "sorry the amount sent is greater than the allowed limit"

		assert.NotNil(test, errorRequest)
		assert.EqualErrorf(test, errorRequest, errorRequest.Error(), "sorry the amount sent is greater than the allowed limit")
		assert.Equal(test, messageExpected, errorRequest.Error())
		assert.Empty(test, responseRequest)
	})
}
