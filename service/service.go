package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/studiously/classsvc/models"
)

type Service interface {
	// GetClasses gets all classes a user is in.
	GetClasses(ctx context.Context) ([]*models.Class, error)
	// GetClass gets details for a specific class.
	GetClass(ctx context.Context, id uuid.UUID) (*models.Class, error)
	// CreateClass creates a class and enrolls the current user in it.
	CreateClass(ctx context.Context, class *models.Class) (uuid.UUID, error)
	// UpdateClass updates a class.
	UpdateClass(ctx context.Context, class *models.Class) error
	// DeleteClass deactivates a class and deletes all associated members.
	DeleteClass(ctx context.Context, id uuid.UUID) error

	// GetMember gets a member by ID.
	GetMember(ctx context.Context, id uuid.UUID) (*models.Member, error)
	//GetMemberByUser(ctx context.Context, user uuid.UUID, class uuid.UUID) (*models.Member, error)
	// CreateMember creates a member (effectively enrolling a user in a class).
	CreateMember(ctx context.Context, member *models.Member) (uuid.UUID, error)
	// DeleteMember deletes a member (effectively un-enrolling a user in a class).
	DeleteMember(ctx context.Context, id uuid.UUID) error
	// UpdateMember updates a member.
	UpdateMember(ctx context.Context, member *models.Member) error
	// GetClassMembers gets all members of a class.
	GetClassMembers(ctx context.Context, classId uuid.UUID) ([]*models.Member, error)
}
