package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fabysdev/fabyscore-go-common/env"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup configures zerolog by using environment variables.
// Errors will exit the process fatally.
func Setup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// set global log level
	levelStr := env.StringDefault("LOG_LEVEL", "info")
	logLevel, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		log.Fatal().Err(err).Str("code", "COMMON-LOG-LEVEL").Str("loglevel", levelStr).Msg("error parsing log level")
	}
	zerolog.SetGlobalLevel(logLevel)

	// create logger
	writers := []io.Writer{}

	if env.BoolDefault("LOG_PRETTY", false) {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false, PartsOrder: []string{zerolog.TimestampFieldName, zerolog.LevelFieldName, zerolog.MessageFieldName, zerolog.CallerFieldName}, FormatCaller: formatCaller})
	} else {
		writers = append(writers, os.Stderr)
	}

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()

	if env.BoolDefault("LOG_CALLER", false) {
		log.Logger = log.Logger.With().Caller().Logger()
	}
}

func formatCaller(i interface{}) string {
	var c string
	if cc, ok := i.(string); ok {
		c = cc
	}
	if len(c) > 0 {
		cwd, err := os.Getwd()
		if err == nil {
			c = strings.TrimPrefix(c, cwd)
			c = strings.TrimPrefix(c, "/")
		}
		c = fmt.Sprintf("\x1b[%dm%v\x1b[0m", 1, c)
	}
	return c
}
