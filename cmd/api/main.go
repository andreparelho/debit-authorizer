package main

import (
	"github.com/andreparelho/debit-authorizer/internal/repository"
	"github.com/andreparelho/debit-authorizer/internal/server"
)

func main() {
	transactions := make(map[string]repository.Client)
	repo := repository.NewClientHistorical(transactions)
	server.RegisterRoutes(repo, transactions)

	server.StartServer()
}
