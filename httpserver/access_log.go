package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// AccessLogMiddleware logs the request using zerolog.
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		rw := &responseWriter{w, 200, 0}

		next.ServeHTTP(rw, r)

		log := log.Info()
		log.Int("b", rw.bytes)
		log.Str("r", r.Referer())
		log.Str("ua", r.UserAgent())
		log.Str("ip", ResolveIP(r))
		log.Int64("d", time.Since(t).Microseconds())
		log.Str("u", r.RequestURI)
		log.Int("s", rw.status)
		log.Str("m", r.Method)
		log.Msg(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
	})
}

// -----------------------------------------------------------------------------------------------------------
type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (rec *responseWriter) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseWriter) Write(b []byte) (int, error) {
	n, err := rec.ResponseWriter.Write(b)

	rec.bytes += n

	return n, err
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
