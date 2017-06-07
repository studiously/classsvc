package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/studiously/classsvc/classsvc"
	"github.com/studiously/introspector"
)

type Middleware func(classsvc.Service) classsvc.Service

func subj(ctx context.Context) uuid.UUID {
	return ctx.Value(introspector.SubjectContextKey).(uuid.UUID)
}
