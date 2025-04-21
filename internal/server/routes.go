package server

import (
	"net/http"

	"github.com/andreparelho/debit-authorizer/internal/handler"
	"github.com/andreparelho/debit-authorizer/internal/repository"
)

const (
	debit_auth_endpoint string = "/authorizer-debit"
)

func RegisterRoutes(repo repository.ClientHistorical, transactions map[string]repository.Client) {
	http.HandleFunc(debit_auth_endpoint, handler.DebitAuthorizerHandler(repo, transactions))
}
