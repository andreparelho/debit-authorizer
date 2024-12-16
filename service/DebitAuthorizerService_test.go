package service_test

import (
	"testing"
	"time"

	request "github.com/andreparelho/debit-authorizer/model/common"
	"github.com/andreparelho/debit-authorizer/service"
	"github.com/stretchr/testify/assert"
)

func TestDebitAuthorizerService(test *testing.T) {

	test.Run("Deve validar com sucesso o debito autorizado", func(test *testing.T) {
		var request = request.RequestAuthorizerDebit{
			ClientId: "1",
			DateTime: time.Now(),
			Amount:   500,
		}

		message, errorService := service.DebitAuthorizerService(request)
		var messageExpected string = `{"message": "debit approved"}`

		assert.Nil(test, errorService)
		assert.Equal(test, messageExpected, string(message))
	})

	test.Run("Deve validar erro quando enviar o segundo request por ser acima do limite permitido", func(test *testing.T) {
		var firstRequest = request.RequestAuthorizerDebit{
			ClientId: "2",
			DateTime: time.Now(),
			Amount:   1000,
		}

		responseFirstRequest, errorFirstRequest := service.DebitAuthorizerService(firstRequest)
		var messageExpected string = `{"message": "debit approved"}`

		assert.Nil(test, errorFirstRequest)
		assert.Equal(test, messageExpected, string(responseFirstRequest))

		var secondRequest = request.RequestAuthorizerDebit{
			ClientId: "2",
			DateTime: time.Now(),
			Amount:   100,
		}

		responseSecondRequest, errorSecondRequest := service.DebitAuthorizerService(secondRequest)
		messageExpected = `{"message": "sorry you have reached your debit limit"}`

		assert.NotNil(test, errorSecondRequest)
		assert.EqualErrorf(test, errorSecondRequest, errorSecondRequest.Error(), "sorry you have reached your debit limit")
		assert.Equal(test, messageExpected, string(responseSecondRequest))
	})

	test.Run("Deve validar erro quando o valor enviado estiver acima do limite maxima por transacao", func(test *testing.T) {
		var request = request.RequestAuthorizerDebit{
			ClientId: "3",
			DateTime: time.Now(),
			Amount:   1500,
		}

		responseRequest, errorRequest := service.DebitAuthorizerService(request)
		var messageExpected string = `{"message": "sorry the amount sent is greater than the allowed limit"}`

		assert.NotNil(test, errorRequest)
		assert.EqualErrorf(test, errorRequest, errorRequest.Error(), "sorry the amount sent is greater than the allowed limit")
		assert.Equal(test, messageExpected, string(responseRequest))
	})
}
