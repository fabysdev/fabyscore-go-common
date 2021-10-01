package gql

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/fabysdev/fabyscore-go/server"
	"github.com/rs/zerolog/log"
)

// AccessLog logs the request using zerolog.
type AccessLog struct {
	// NoOperationNameError if true a requst with no operation name will be logged as a warning
	NoOperationNameError bool
}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
} = AccessLog{}

// ExtensionName is the name of the extension.
func (AccessLog) ExtensionName() string {
	return "AccessLog"
}

// Validate is called when adding an extension to the server, it allows validation against the servers schema.
func (AccessLog) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptResponse is the interceptor to log the request.
func (al AccessLog) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	rc := graphql.GetOperationContext(ctx)

	start := rc.Stats.OperationStart
	resp := next(ctx)
	end := graphql.Now()

	errors := graphql.GetErrors(ctx)
	requestInfos := ctx.Value(requestInformationContextKey).(*requestInformation)

	loglevel := zerolog.InfoLevel

	var errorCodes []string
	if errors != nil {
		loglevel = zerolog.WarnLevel

		errorCodes = make([]string, len(errors))
		for i, err := range errors {
			if code, ok := err.Extensions["code"]; ok {
				if code == "PERSISTED_QUERY_NOT_FOUND" {
					return resp
				}

				errorCodes[i] = code.(string)
			}
		}
	}

	opName := rc.OperationName
	if al.NoOperationNameError && opName == "" {
		opName = "[NO OPERATION NAME] " + rc.RawQuery
		loglevel = zerolog.WarnLevel
	}

	l := log.WithLevel(loglevel)
	l.Str("r", requestInfos.r)
	l.Str("ua", requestInfos.ua)
	l.Str("ip", requestInfos.ip)
	l.Int64("d", end.Sub(start).Microseconds())

	if errors != nil {
		l.Err(errors)
		l.Strs("codes", errorCodes)
		l.Msg(fmt.Sprintf("error: %s", opName))
	} else {
		l.Msg(fmt.Sprintf("%s: %s", string(rc.Operation.Operation), opName))
	}

	return resp
}

// -----------------------------------------------------------------------------------------------------------
type requestInformation struct {
	r  string
	ua string
	ip string
}

var requestInformationContextKey = &server.ContextKey{Name: "requestinformation"}

// AccessLogMiddleware adds request information to the context.
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := &requestInformation{
			r:  r.Referer(),
			ua: r.UserAgent(),
			ip: ResolveIP(r),
		}

		ctx := context.WithValue(r.Context(), requestInformationContextKey, info)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// -----------------------------------------------------------------------------------------------------------
var xForwardedForHeaderKey = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIPHeaderKey = http.CanonicalHeaderKey("X-Real-IP")

// ResolveIP returns the client ip.
func ResolveIP(r *http.Request) string {
	// X-Forwarded-For
	ip := r.Header.Get(xForwardedForHeaderKey)
	if ip != "" {
		i := strings.Index(ip, ",")
		if i != -1 {
			ip = ip[:i]
		}

		return ip
	}

	// X-Real-IP
	ip = r.Header.Get(xRealIPHeaderKey)
	if ip != "" {
		return ip
	}

	// RemoteAddr
	return r.RemoteAddr
}
