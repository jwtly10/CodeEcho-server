package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

		environment := os.Getenv("GO_ENV")
		if environment == "" {
			// By default, assume production, so we dont log sensitive information
			environment = "prod"
		}

		if environment == "dev" {
			output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
			log = zerolog.New(output).
				Level(zerolog.DebugLevel).
				With().
				Timestamp().
				Caller().
				Int("pid", os.Getpid()).
				Logger()
		} else {
			// TODO - Add a log file for production
			output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
			log = zerolog.New(output).
				Level(zerolog.InfoLevel).
				With().
				Timestamp().
				Logger()
		}
	})
	return log
}
