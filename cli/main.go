package main

import (
	"fmt"

	request "github.com/andreparelho/debit-authorizer/model/common"
	"github.com/andreparelho/debit-authorizer/service"
)

func main() {
	var endStatus = true

	for endStatus != false {
		var clientId string
		var amount float64
		var endString string

		fmt.Print("Cliente ID: ")
		fmt.Scanln(&clientId)
		fmt.Print("Amount: ")
		fmt.Scanln(&amount)

		var request = request.RequestAuthorizerDebit{
			ClientId: clientId,
			Amount:   amount,
		}

		service.DebitAuthorizerService(request)

		fmt.Print("Deseja finalizar? SIM / NAO: ")
		if endString == "SIM" {
			endStatus = false
		}
	}
}
