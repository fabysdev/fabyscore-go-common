package gql

import (
	"context"
	"errors"
	"strings"

	"github.com/fabysdev/fabyscore-go-common/env"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateServer returns a new graphql server with middlewares.
func CreateServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	if env.BoolDefault("GRAPHQL_INTROSPECTION", false) {
		srv.Use(extension.Introspection{})
	}

	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	if env.BoolDefault("GRAPHQL_TRACING", false) {
		srv.Use(apollotracing.Tracer{})
	}

	srv.Use(AccessLog{NoOperationNameError: env.BoolDefault("LOG_OPNAMERROR", true)})

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		code, ok := err.Extensions["code"]
		if !ok {
			log.Error().Err(err).Str("code", "ERROR-UNKNOWN").Msg(err.Error())

			if err.Extensions == nil {
				err.Extensions = map[string]interface{}{}
			}
			err.Extensions["code"] = "ERROR-UNKNOWN"
		} else if strings.HasPrefix(code.(string), "GRAPHQL") {
			log.Error().Err(err).Interface("code", code).Msg(err.Error())
		}

		return err
	})

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Error().Interface("err", err).Str("code", "ERROR-INTERNAL").Msg("graphql internal server panic error")
		return errors.New("unkown error")
	})

	return srv
}
