package service

import (
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/studiously/classsvc/models"
)

var (
	ErrUnauthenticated = errors.New("must be authenticated")
)

type Service interface {
	GetClass(id uuid.UUID) (*models.Class, error)
	CreateClass(class *models.Class) (error)
	UpdateClass(class *models.Class) error
	DeleteClass(id uuid.UUID) error

	GetMember(id uuid.UUID) (*models.Member, error)
	GetMemberByUser(user uuid.UUID, class uuid.UUID) (*models.Member, error)
	CreateMember(member *models.Member) error
	DeleteMember(id uuid.UUID) error
	UpdateMember(member *models.Member) error
	GetClassMembers(classId uuid.UUID) ([]*models.Member, error)
}
