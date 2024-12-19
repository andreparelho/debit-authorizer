package repository_test

import (
	"testing"
	"time"

	model "github.com/andreparelho/debit-authorizer/model/service"
	"github.com/andreparelho/debit-authorizer/repository"
	"github.com/stretchr/testify/assert"
)

func TestClientTransactionHistoricalRepository(test *testing.T) {
	test.Run("Deve validar a criacao do client e o historico quando chamar a funcao CreateClientHistorical", func(test *testing.T) {
		var transactionHistorical = make(map[string]model.Client)
		var dateTime time.Time = time.Now()
		var amount float64 = 100

		var client model.Client = model.Client{
			ClientId:    "Andre",
			LastPayment: dateTime,
			TotalAmount: amount,
			Historical:  []model.Historical{},
		}

		assert.Empty(test, client.Historical)

		repository.CreateClientHistorical(transactionHistorical, client, dateTime, amount)

		client = transactionHistorical[client.ClientId]

		historical := client.Historical

		assert.NotEmpty(test, client)
		assert.NotEmpty(test, client.Historical)
		assert.True(test, client.ClientId == "Andre")
		assert.True(test, client.LastPayment == dateTime)
		assert.True(test, client.TotalAmount == amount)
		for _, historic := range historical {
			assert.True(test, historic.Amount == amount)
			assert.True(test, historic.DateTime == dateTime)
		}
	})

	test.Run("Deve validar a atualzacao do historico quando chamar a funcao UpdateClientHistorical", func(test *testing.T) {
		var transactionHistorical = make(map[string]model.Client)
		var dateTime time.Time = time.Now()

		var client model.Client = model.Client{
			ClientId:    "1",
			LastPayment: dateTime,
			TotalAmount: 100,
			Historical:  []model.Historical{},
		}

		repository.CreateClientHistorical(transactionHistorical, client, dateTime, 100)
		assert.True(test, client.TotalAmount == 100)

		client = transactionHistorical[client.ClientId]

		client = model.Client{
			ClientId:    client.ClientId,
			LastPayment: client.LastPayment,
			TotalAmount: client.TotalAmount,
			Historical:  client.Historical,
		}

		var amount float64 = 500
		var totalAmount float64 = client.TotalAmount + amount
		repository.UpdateClientHistorical(client, transactionHistorical, client.ClientId, dateTime, totalAmount, amount)

		client = transactionHistorical[client.ClientId]

		assert.True(test, client.ClientId == "1")
		assert.True(test, client.TotalAmount == 600)
		assert.True(test, len(client.Historical) == 2)

		historical := client.Historical

		assert.True(test, historical[0].Amount == 100)
		assert.True(test, historical[1].Amount == 500)
	})

	test.Run("Deve retornar um client valido quando chamar GetClientHitorical", func(test *testing.T) {
		var transactionHistorical = make(map[string]model.Client)
		var dateTime time.Time = time.Now()
		var amount float64 = 100

		var client model.Client = model.Client{
			ClientId:    "Andre",
			LastPayment: dateTime,
			TotalAmount: amount,
			Historical:  []model.Historical{},
		}

		assert.Empty(test, client.Historical)

		repository.CreateClientHistorical(transactionHistorical, client, dateTime, amount)

		client = transactionHistorical[client.ClientId]

		var response model.Client = repository.GetClientHitorical(client.ClientId, transactionHistorical)

		assert.NotEmpty(test, response)
		assert.Equal(test, response, client)
		assert.Equal(test, response.ClientId, client.ClientId)
		assert.Equal(test, response.Historical, client.Historical)
		assert.Equal(test, response.TotalAmount, client.TotalAmount)
	})
}
