package logger

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/pkgerrors"

	"github.com/michaljirman/goblockchain/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setups up a global logger based on a logging configuration provided.
func SetupLogger(cfg *config.LogConf) error {
	zerolog.TimeFieldFormat = ""
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Error().Err(err).Msg("error returned when parsing log level")
		return err
	}
	zerolog.SetGlobalLevel(level)

	if cfg.DevelopmentLogger {
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			fmt.Println(string(debug.Stack()))
			return nil
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Stack().Logger()
	} else {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Stack().Logger()

	}
	return nil
}
