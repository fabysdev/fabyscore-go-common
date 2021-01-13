package gql

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Error returns a gqlerror Error with the code extension set. Every error is logged.
func Error(errorCode string, msg string, level zerolog.Level, err error, logMsg string) *gqlerror.Error {
	msgToLog := msg
	if logMsg != "" {
		msgToLog = logMsg
	}

	log.WithLevel(level).Err(err).Str("code", errorCode).Msg(msgToLog)

	return &gqlerror.Error{
		Message:    msg,
		Extensions: map[string]interface{}{"code": errorCode},
	}
}

// Err returns an actual error e.g. for db errors.
func Err(errorCode string, msg string, err error, logMsg string) *gqlerror.Error {
	return Error(errorCode, msg, zerolog.ErrorLevel, err, logMsg)
}

// ErrBadRequest returns a standard error if something in the request is not valid.
func ErrBadRequest(errorCode string, msg string, err error) *gqlerror.Error {
	level := zerolog.InfoLevel
	if err != nil {
		level = zerolog.WarnLevel
	}

	return Error(errorCode, msg, level, err, "")
}
