package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ory/hydra/oauth2"
	"github.com/studiously/classsvc/middleware/auth"
	"github.com/studiously/classsvc/models"
)

type postgresService struct {
	*sql.DB
}

func (s *postgresService) GetClass(ctx context.Context, classId uuid.UUID) (*models.Class, error) {
	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
	subj, err := uuid.Parse(introspection.Subject)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	_, err = models.MemberByUserIDClassID(s, subj, classId)
	if err != nil {
		// Return ErrNotFound to protect the secrecy of the class (whether or not it exists)
		return nil, ErrNotFound
	}
	return models.ClassByID(s, classId)
}

func (s *postgresService) CreateClass(ctx context.Context, class *models.Class) (uuid.UUID, error) {
	class.ID = uuid.New()
	return class.ID, class.Save(s)
}

func (s *postgresService) UpdateClass(ctx context.Context, class *models.Class) error {
	return class.Update(s)
}

func (s *postgresService) DeleteClass(ctx context.Context, classId uuid.UUID) error {
	class, err := models.ClassByID(s, classId)
	if err != nil {
		return err
	}
	err = class.Delete(s)
	return err
}

func (s *postgresService) GetClassMembers(ctx context.Context, classId uuid.UUID) ([]*models.Member, error) {
	return models.MembersByClassID(s, classId)
}

func (s *postgresService) GetMember(ctx context.Context, id uuid.UUID) (*models.Member, error) {
	return models.MemberByID(s, id)
}

func (s *postgresService) GetMemberByUser(ctx context.Context, userId uuid.UUID, classId uuid.UUID) (*models.Member, error) {
	return models.MemberByUserIDClassID(s, userId, classId)
}

func (s *postgresService) CreateMember(ctx context.Context, member *models.Member) (uuid.UUID, error) {
	member.ID = uuid.New()
	return member.ID, member.Save(s)
}

func (s *postgresService) DeleteMember(ctx context.Context, id uuid.UUID) error {
	member, err := models.MemberByID(s, id)
	if err != nil {
		return err
	}
	err = member.Delete(s)
	return err
}

func (s *postgresService) UpdateMember(ctx context.Context, member *models.Member) error {
	return member.Update(s)
}

func (s *postgresService) GetClasses(ctx context.Context) ([]*models.Class, error) {
	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
	subj, err := uuid.Parse(introspection.Subject)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	members, err := models.MembersByUserID(s, subj)
	if err != nil {
		return nil, err
	}
	var results []*models.Class
	for i := range members {
		results[i], err = models.ClassByID(s, members[i].ClassID)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func NewPostgres(db *sql.DB) Service {
	return &postgresService{db}
}
