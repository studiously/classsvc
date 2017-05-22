package service

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/studiously/classsvc/models"
)

var (
	ErrUnauthenticated = errors.New("either no token was passed or the provided token was malformed")
	ErrForbidden = errors.New("requester does not have correct scope for resource")
	ErrNotFound = errors.New("not found")
)

type Service interface {
	GetClasses(ctx context.Context) ([]*models.Class, error)
	GetClass(ctx context.Context, id uuid.UUID) (*models.Class, error)
	CreateClass(ctx context.Context, class *models.Class) (uuid.UUID, error)
	UpdateClass(ctx context.Context, class *models.Class) error
	DeleteClass(ctx context.Context, id uuid.UUID) error

	GetMember(ctx context.Context, id uuid.UUID) (*models.Member, error)
	GetMemberByUser(ctx context.Context, user uuid.UUID, class uuid.UUID) (*models.Member, error)
	CreateMember(ctx context.Context, member *models.Member) (uuid.UUID, error)
	DeleteMember(ctx context.Context, id uuid.UUID) error
	UpdateMember(ctx context.Context, member *models.Member) error
	GetClassMembers(ctx context.Context, classId uuid.UUID) ([]*models.Member, error)
}
