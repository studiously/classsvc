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

func (s *postgresService) CreateClass(ctx context.Context, name string) (uuid.UUID, error) {
	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
	subj, err := uuid.Parse(introspection.Subject)
	if err != nil {
		return uuid.Nil, ErrUnauthenticated
	}
	tx, err := s.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		return uuid.Nil, err
	}
	class := models.Class{
		ID:          uuid.New(),
		Name:        name,
		CurrentUnit: uuid.Nil,
		Active:      true,
	}
	err = class.Save(tx)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	member := models.Member{
		ID:      uuid.New(),
		UserID:  subj,
		ClassID: class.ID,
		Role:    models.UserRoleAdministrator,
	}
	err = member.Save(tx)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	err = tx.Commit()
	return class.ID, err
}

func (s *postgresService) UpdateClass(ctx context.Context, id uuid.UUID, name string, currentUnit uuid.UUID) error {
	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
	subj, err := uuid.Parse(introspection.Subject)
	if err != nil {
		return ErrUnauthenticated
	}
	member, err := models.MemberByUserIDClassID(s, subj, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrForbidden
		default:
			return err
		}
	}
	if member.Role < models.UserRoleAdministrator {
		return ErrForbidden
	}
	class, err := models.ClassByID(s, id)
	if err != nil {
		return err
	}
	if name != "" {
		class.Name = name
	}
	if currentUnit != uuid.Nil {
		class.CurrentUnit = currentUnit // TODO need to validate current unit.
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
	if member.Role != models.UserRoleAdministrator {
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

func (s *postgresService) ListMembers(ctx context.Context, class uuid.UUID) ([]*models.Member, error) {
	_, err := models.MemberByUserIDClassID(s, subj(ctx), class)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Either user is not in class or class does not exist.
			return nil, ErrForbidden
		default:
			return nil, err
		}
	}
	return models.MembersByClassID(s, class)
}

//func (s *postgresService) GetMember(ctx context.Context, id uuid.UUID) (*models.Member, error) {
//	rm, err := models.MemberByID(s, id)
//	if err != nil {
//		switch err {
//		case sql.ErrNoRows:
//			return nil, ErrNotFound
//		default:
//			return nil, err
//		}
//	}
//	m, err := models.MemberByUserIDClassID(s, subj(ctx), rm.ClassID)
//	if err != nil {
//		switch err {
//		case sql.ErrNoRows:
//			return nil, ErrForbidden
//		default:
//			return nil, err
//		}
//	}
//	return m, nil
//}

//func (s *postgresService) GetMemberByUser(ctx context.Context, userId uuid.UUID, classId uuid.UUID) (*models.Member, error) {
//	introspection := ctx.Value(auth.OAuth2IntrospectionContextKey).(oauth2.Introspection)
//	subj, err := uuid.Parse(introspection.Subject)
//	if err != nil {
//		return nil, ErrUnauthenticated
//	}
//	if userId !=
//	return models.MemberByUserIDClassID(s, userId, classId)
//}

func (s *postgresService) JoinClass(ctx context.Context, class uuid.UUID) error {
	_, err := models.MemberByUserIDClassID(s, subj(ctx), class)
	if err == nil {
		return ErrUserEnrolled
	}
	member := models.Member{
		ID:      uuid.New(),
		UserID:  subj(ctx),
		ClassID: class,
		Role:    models.UserRoleStudent,
	}
	return member.Save(s)
}

func (s *postgresService) LeaveClass(ctx context.Context, user uuid.UUID, class uuid.UUID) error {
	self, err := models.MemberByUserIDClassID(s, subj(ctx), class)
	if err != nil {
		return err
	}
	if user != subj(ctx) {
		if self.Role != models.UserRoleAdministrator {
			return ErrForbidden
		}
		target, err := models.MemberByUserIDClassID(s, user, class)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return ErrNotFound
			default:
				return err
			}
		}
		if target.Role >= self.Role {
			// Can't remove equal or superior.
			return ErrForbidden
		}
		return target.Delete(s)
	} else {
		if self.Role == models.UserRoleOwner {
			return ErrMustSetOwner
		}
		return self.Delete(s)
	}
}

func (s *postgresService) SetRole(ctx context.Context, user uuid.UUID, class uuid.UUID, role models.UserRole) error {
	// Can't set own role directly
	if user == subj(ctx) {
		return ErrForbidden
	}
	self, err := models.MemberByUserIDClassID(s, subj(ctx), class)
	if err != nil {
		return err
	}
	target, err := models.MemberByUserIDClassID(s, user, class)
	if err != nil {
		return err
	}
	// Can't set role of users superior or equal.
	if self.Role <= target.Role {
		return ErrForbidden
	}
	// Special case: setting another owner means making the current owner an administrator.
	if self.Role == models.UserRoleOwner && role == models.UserRoleOwner {
		tx, err := s.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
		if err != nil {
			return err
		}
		self.Role = models.UserRoleAdministrator
		err = self.Save(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		target.Role = models.UserRoleOwner
		err = target.Save(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
	target.Role = role
	return target.Save(s)
}
func (s *postgresService) ListClasses(ctx context.Context) ([]uuid.UUID, error) {
	members, err := models.MembersByUserID(s, subj(ctx))
	if err != nil {
		return nil, err
	}
	var results []uuid.UUID
	for i := range members {
		results[i] = members[i].ClassID
	}
	return results, nil
}

func NewPostgres(db *sql.DB) Service {
	return &postgresService{db}
}

func subj(ctx context.Context) uuid.UUID {
	return ctx.Value(auth.SubjectContextKey).(uuid.UUID)
}
