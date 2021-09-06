package logger

import (
	"cloud.google.com/go/logging"
	"github.com/rs/zerolog"
)

var severityLevels = map[zerolog.Level]string{
	zerolog.DebugLevel: logging.Debug.String(),
	zerolog.InfoLevel:  logging.Info.String(),
	zerolog.WarnLevel:  logging.Warning.String(),
	zerolog.ErrorLevel: logging.Error.String(),
	zerolog.FatalLevel: logging.Critical.String(),
	zerolog.PanicLevel: logging.Emergency.String(),
}

// SeverityHook adds the "severity" field for google cloud structured logging.
type SeverityHook struct{}

// Run adds the "severity" field to the given log event.
func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		if severity, ok := severityLevels[level]; ok {
			e.Str("severity", severity)
		}
	}
}
