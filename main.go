package Main

import (
	"fmt"
	"net/http"

	"github.com/andreparelho/debit-authorizer/handler"
)

func main() {
	http.HandleFunc("/approver", handler.ApproverHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
