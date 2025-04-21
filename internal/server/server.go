package server

import (
	"log"
	"net/http"
)

func StartServer() {
	log.Println("Started on port", ":8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
