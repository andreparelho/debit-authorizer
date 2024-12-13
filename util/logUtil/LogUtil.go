package util

import (
	"os"
	"strconv"

	serviceDTO "github.com/andreparelho/credit-approver/model/service/dto"
	"github.com/rs/zerolog"
)

func ServiceLoggerInfo(client serviceDTO.Client, clientId string, message string) {
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	logger.Info().
		Str("client", clientId).
		Str("lastPayment", client.LastPayment.String()).
		Str("totalAmount", strconv.Itoa(int(client.TotalAmount))).
		Msg(message)
}
