package service

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ory/hydra/oauth2"

)

var (
	TestUserId = uuid.New()
)

func GetClass(t *testing.T, s Service) {
	c := newClass()
	s.CreateClass(ContextWithClaims(), c)
}

func ContextWithIntrospection(ctx context.Context) context.Context {

}

type claims struct {
	jwt.StandardClaims
	scopes oauth2.Introspection
}

func newClass() Class {
	return Class{
		Id:      uuid.New(),
		Members: make([]uuid.UUID, 0),
	}
}
