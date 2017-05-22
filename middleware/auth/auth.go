package auth

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	ajwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	hoauth2 "github.com/ory/hydra/oauth2"
)

type contextKey string

const (
	// OAuth2TokenContextKey holds the key used to store an OAuth2 Token in the context.
	OAuth2IntrospectionContextKey contextKey = "OAuth2Introspection"
)

func New(introspector hoauth2.Introspector, scopes ...string) endpoint.Middleware {
	return func(outer endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			introspection, err := introspector.IntrospectToken(ctx, ctx.Value(ajwt.JWTTokenContextKey).(jwt.Token).Raw, scopes...)
			if err != nil {
				return nil, err
			}
			return outer(context.WithValue(ctx, OAuth2IntrospectionContextKey, introspection), request)
		}
	}
}