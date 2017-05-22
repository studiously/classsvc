package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ory/hydra/oauth2"
	. "github.com/studiously/classsvc/errors"
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
	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
	subj, err := uuid.Parse(introspection.Subject)
	if err != nil {
		return uuid.Nil, ErrUnauthenticated
	}
	tx, err := s.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		return uuid.Nil, err
	}
	class.ID = uuid.New()
	class.CurrentUnit = uuid.Nil
	err = class.Save(tx)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	member := models.Member{
		ID:      uuid.New(),
		UserID:  subj,
		ClassID: class.ID,
		Role:    models.MemberRoleAdministrator,
	}
	err = member.Save(tx)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	err = tx.Commit()
	return class.ID, err
}

func (s *postgresService) UpdateClass(ctx context.Context, class *models.Class) error {
	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
	subj, err := uuid.Parse(introspection.Subject)
	if err != nil {
		return ErrUnauthenticated
	}
	member, err := models.MemberByUserIDClassID(s, subj, class.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrForbidden
		default:
			return err
		}
	}
	if member.Role != models.MemberRoleAdministrator {
		return ErrForbidden
	}
	return class.Update(s)
}

func (s *postgresService) DeleteClass(ctx context.Context, classId uuid.UUID) error {
	member, err := models.MemberByUserIDClassID(s, subj(ctx), classId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Either user is not in class or class does not exist.
			return ErrForbidden
		default:
			return err
		}
	}
	if member.Role != models.MemberRoleAdministrator {
		return ErrForbidden
	}
	class, err := models.ClassByID(s, classId)
	if err != nil {
		return err
	}
	tx, err := s.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		tx.Rollback()
		return err
	}
	class.Active = false
	class.Save(tx)
	tx.Exec("DELETE FROM members WHERE class_id=$1;", class.ID)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func (s *postgresService) GetClassMembers(ctx context.Context, classId uuid.UUID) ([]*models.Member, error) {
	_, err := models.MemberByUserIDClassID(s, subj(ctx), classId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Either user is not in class or class does not exist.
			return nil, ErrForbidden
		default:
			return nil, err
		}
	}
	return models.MembersByClassID(s, classId)
}

func (s *postgresService) GetMember(ctx context.Context, id uuid.UUID) (*models.Member, error) {
	rm, err := models.MemberByID(s, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	m, err := models.MemberByUserIDClassID(s, subj(ctx), rm.ClassID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrForbidden
		default:
			return nil, err
		}
	}
	return m, nil
}

//func (s *postgresService) GetMemberByUser(ctx context.Context, userId uuid.UUID, classId uuid.UUID) (*models.Member, error) {
//	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
//	subj, err := uuid.Parse(introspection.Subject)
//	if err != nil {
//		return nil, ErrUnauthenticated
//	}
//	if userId !=
//	return models.MemberByUserIDClassID(s, userId, classId)
//}

func (s *postgresService) CreateMember(ctx context.Context, member *models.Member) (uuid.UUID, error) {
	member.ID = uuid.New()
	member.UserID = subj(ctx)
	return member.ID, member.Save(s)
}

func (s *postgresService) DeleteMember(ctx context.Context, id uuid.UUID) error {
	member, err := models.MemberByID(s, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}
	if member.UserID.String() != subj(ctx).String() {
		return ErrForbidden
	}

	err = member.Delete(s)
	return err
}

func (s *postgresService) UpdateMember(ctx context.Context, member *models.Member) error {
	if member.ID.String() != subj(ctx).String() {
		return ErrForbidden
	}
	return member.Update(s)
}

func (s *postgresService) GetClasses(ctx context.Context) ([]*models.Class, error) {
	members, err := models.MembersByUserID(s, subj(ctx))
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

func subj(ctx context.Context) uuid.UUID {
	return ctx.Value(auth.SubjectContextKey).(uuid.UUID)
}
