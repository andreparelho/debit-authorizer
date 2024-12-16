package main

import (
	"fmt"
	"strings"
	"time"

	request "github.com/andreparelho/debit-authorizer/model/common"
	"github.com/andreparelho/debit-authorizer/service"
)

const END_STRING string = "sim"

func main() {
	var endStatus bool = true

	for endStatus {
		var clientId string
		var amount float64
		var endString string
		var now time.Time = time.Now()

		fmt.Print("Client ID: ")
		fmt.Scanln(&clientId)
		fmt.Print("Amount: ")
		fmt.Scanln(&amount)

		var request = request.RequestAuthorizerDebit{
			ClientId: clientId,
			DateTime: now,
			Amount:   amount,
		}

		response, errorService := service.DebitAuthorizerService(request)
		if errorService != nil {
			fmt.Println(errorService.Error())
		} else {
			fmt.Println(response)
		}

		fmt.Print("Deseja finalizar? SIM / NAO: ")
		fmt.Scanln(&endString)
		if strings.Compare(strings.ToLower(endString), strings.ToLower(END_STRING)) == 0 {
			endStatus = false
		}
	}
}
