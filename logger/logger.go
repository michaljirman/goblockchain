package logger

import (
	"os"

	"github.com/rs/zerolog/pkgerrors"

	"github.com/michaljirman/goblockchain/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(cfg *config.LogConf) error {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = ""

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Error().Err(err).Msg("error returned when parsing log level")
		return err
	}
	zerolog.SetGlobalLevel(level)

	if cfg.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Caller().Logger()
	} else {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}
	return nil
}
