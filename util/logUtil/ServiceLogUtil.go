package util

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

func ServiceLoggerInfo(lastPayment time.Time, totalAmount float64, clientId string, message string) {
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	logger.Info().
		Str("class", "SERVICE").
		Str("client", clientId).
		Str("lastPayment", lastPayment.String()).
		Str("totalAmount", strconv.Itoa(int(totalAmount))).
		Msg(message)
}

func ServiceLoggerError(lastPayment time.Time, totalAmount float64, clientId string, message string) {
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	logger.Error().
		Str("class", "SERVICE").
		Str("client", clientId).
		Str("lastPayment", lastPayment.String()).
		Str("totalAmount", strconv.Itoa(int(totalAmount))).
		Msg(message)
}
