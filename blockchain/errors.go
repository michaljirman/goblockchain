package blockchain

import "github.com/rs/zerolog/log"

func HandleError(err error) {
	if err != nil {
		log.Error().Stack().Err(err).Msg("unexpected error")
	}
}

func HandleFatal(err error) {
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("unexpected error")
	}
}
