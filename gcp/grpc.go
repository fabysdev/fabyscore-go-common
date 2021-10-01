package gcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceToken struct {
	tokenSource oauth2.TokenSource
}

// newServiceToken returns the credentials.PerRPCCredentials to add the auth token for the given grpcAddr.
// grpcAddr must be of form "addr:port" without https, e.g. "localhost:8080"
func newServiceToken(grpcAddr string) credentials.PerRPCCredentials {
	i := strings.LastIndexByte(grpcAddr, ':')
	serviceURL := "https://" + grpcAddr[:i]

	tokenSource, err := idtoken.NewTokenSource(context.Background(), serviceURL)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create auth token source")
	}

	return serviceToken{
		tokenSource: tokenSource,
	}
}

// GetRequestMetadata adds the auth token to the request metadata.
func (t serviceToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := t.tokenSource.Token()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve service auth token")
	}

	return map[string]string{
		"authorization": "Bearer " + token.AccessToken,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (serviceToken) RequireTransportSecurity() bool {
	return true
}

// NewGRPCConn creates a new grpc client connection with service token authentication.
// host must be of form "addr:port" without a scheme, e.g. "localhost:8080"
// no authentication token is added if insecure is true.
func NewGRPCConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithAuthority(host))

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))

		opts = append(opts, grpc.WithPerRPCCredentials(newServiceToken(host)))
	}

	return grpc.Dial(host, opts...)
}
