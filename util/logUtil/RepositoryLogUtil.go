package util

import (
	"encoding/json"
	"os"

	model "github.com/andreparelho/debit-authorizer/model/service"
	"github.com/rs/zerolog"
)

func RepositoryLoggerInfo(clientId string, historical []model.Historical, message string) {
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	historicalJSON, _ := json.Marshal(historical)

	logger.Info().
		Str("class", "REPOSITORY").
		Str("client", clientId).
		Str("historical", string(historicalJSON)).
		Msg(message)
}
