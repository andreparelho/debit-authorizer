package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/andreparelho/debit-authorizer/internal/repository"
	"github.com/rs/zerolog"
)

const (
	LAST_FIVE_MINUTES = 5 * time.Minute
	MAX_TOTAL_AMOUNT  = 1000
)

var mutex sync.Mutex

func DebitAuthorizerHandler(repo repository.ClientHistorical, transactions map[string]repository.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

		if r.Method != http.MethodPost {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", "wrong method")

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", "error to read body request")

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var request RequestAuthorizerDebit
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", "error to unmarshal request")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		errValidate := validateRequest(request)
		if errValidate != nil {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", errValidate.Error())

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		now := time.Now()
		dateTime := getDate(request.DateTime, now)

		mutex.Lock()
		var clientId = request.ClientId
		client, ok := transactions[clientId]
		mutex.Unlock()

		totalAmount := client.TotalAmount + request.Amount
		if totalAmount > MAX_TOTAL_AMOUNT && now.Sub(client.LastPayment) <= LAST_FIVE_MINUTES {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", "sorry you have reached your debit limit")

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if totalAmount > MAX_TOTAL_AMOUNT {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", "sorry the amount sent is greater than the allowed limit")

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		validateClient(ok, client, clientId, dateTime, request.Amount, totalAmount, repo)

		responseRepository := repo.GetClientHistorical(clientId)
		response := ResponseAuthorizerDebit{
			Message:          "debit authorized",
			ClientHistorical: responseRepository,
		}

		var responseHandler []byte
		if responseHandler, err = json.Marshal(response); err != nil {
			logger.Error().
				Str("component", "handler.DebitAuthorizerHandler").
				Str("erro", "error to marshal response repository")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info().
			Str("id", responseRepository.ClientId).
			Time("last_payment", responseRepository.LastPayment).
			Int("historics", len(responseRepository.Historical)).
			Float64("total_amount", responseRepository.TotalAmount)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseHandler)
	}
}

func validateRequest(request RequestAuthorizerDebit) error {
	switch {
	case request.ClientId == "":
		return errors.New("campo id vazio")
	case request.Amount < 0.01:
		return errors.New("amount vazio ou menor que 0.01")
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
