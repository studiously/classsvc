package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	hoauth2 "github.com/ory/hydra/oauth2"
	. "github.com/studiously/classsvc/errors"
)

type contextKey string

const (
	// OAuth2TokenContextKey holds the key used to store an OAuth2 Token in the context.
	OAuth2IntrospectionContextKey contextKey = "OAuth2Introspection"
	SubjectContextKey             contextKey = "Subject"
	TokenContextKey               contextKey = "Token"
)

func New(introspector hoauth2.Introspector, scopes ...string) endpoint.Middleware {
	return func(outer endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			introspection, err := introspector.IntrospectToken(ctx, ctx.Value(TokenContextKey).(string), scopes...)
			if err != nil {
				return nil, err
			}
			subj, err := uuid.Parse(introspection.Subject)
			if err != nil {
				return nil, ErrUnauthenticated
			}
			return outer(context.WithValue(context.WithValue(ctx, OAuth2IntrospectionContextKey, introspection), SubjectContextKey, subj), request)
		}
	}
}
