package main

import (
	"fmt"
	"net/http"

	"github.com/andreparelho/debit-authorizer/api/handler"
)

const ENDPOINT_AUTHORIZER_DEBIT string = "/authorizer-debit"
const SERVER_PORT string = ":8080"

func main() {
	http.HandleFunc(ENDPOINT_AUTHORIZER_DEBIT, handler.DebitAuthorizerHandler)

	fmt.Println("Application is running on port 8080")
	http.ListenAndServe(SERVER_PORT, nil)
}
