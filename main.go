package main

import (
	"fmt"
	"net/http"

	"github.com/andreparelho/credit-approver/handler"
)

func main() {
	http.HandleFunc("/approver", handler.ApproverHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
