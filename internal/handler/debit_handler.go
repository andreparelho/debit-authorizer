package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/andreparelho/debit-authorizer/internal/repository"
)

const (
	LAST_FIVE_MINUTES = 5 * time.Minute
	MAX_TOTAL_AMOUNT  = 1000
)

var mutex sync.Mutex

func DebitAuthorizerHandler(repo repository.ClientHistorical, transactions map[string]repository.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var request RequestAuthorizerDebit
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := validateRequest(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var now time.Time = time.Now()
		var dateTime time.Time = getDate(request.DateTime, now)

		mutex.Lock()
		var clientId = request.ClientId
		client, ok := transactions[clientId]
		mutex.Unlock()

		var totalAmount = client.TotalAmount + request.Amount
		if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if totalAmount > MAX_TOTAL_AMOUNT {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		validateClient(ok, client, clientId, dateTime, request.Amount, totalAmount, repo)

		responseRepository := repo.GetClientHitorical(clientId)
		var response []byte
		if response, err = json.Marshal(responseRepository); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func validateRequest(request RequestAuthorizerDebit) error {
	switch {
	case request.ClientId == "":
		return errors.New("propertie clientId is empty")
	case request.Amount < 0.01:
		return errors.New("propertie amount is empty or less than minimum")
	default:
		return nil
	}
}

func getDate(requestDate time.Time, now time.Time) time.Time {
	if requestDate.IsZero() {
		return now
	}

	return requestDate
}

func validateClient(ok bool, client repository.Client, clientId string, dateTime time.Time, amount float64, totalAmount float64, repo repository.ClientHistorical) {
	if !ok {
		client = repository.Client{
			ClientId:    clientId,
			LastPayment: dateTime,
			TotalAmount: amount,
			Historical:  []repository.Historical{},
		}
		repo.CreateClientHistorical(client, dateTime, amount)
	} else {
		repo.UpdateClientHistorical(client, clientId, dateTime, totalAmount, amount)
	}
}
