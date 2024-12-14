package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andreparelho/debit-authorizer/handler"
	"github.com/joho/godotenv"
)

const ENDPOINT_AUTHORIZER_DEBIT string = "/authorizer-debit"

func main() {
	var errorEnv error = godotenv.Load()
	if errorEnv != nil {
		log.Fatalf("Error to load file .env: %v", errorEnv)
	}

	var serverPort string = os.Getenv("SERVER_PORT")

	http.HandleFunc(ENDPOINT_AUTHORIZER_DEBIT, handler.DebitAuthorizerHandler)

	fmt.Println("Application is running on port 8080")
	http.ListenAndServe(serverPort, nil)
}
