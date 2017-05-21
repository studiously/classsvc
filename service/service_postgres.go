package service

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/studiously/classsvc/models"
)

type postgresService struct {
	*sql.DB
}

func (s *postgresService) GetClass(classId uuid.UUID) (*models.Class, error) {
	return models.ClassByID(s, classId)
}

func (s *postgresService) CreateClass(class *models.Class) (error) {
	return class.Save(s)
}

func (s *postgresService) UpdateClass(class *models.Class) error {
	return class.Update(s)
}

func (s *postgresService) DeleteClass(classId uuid.UUID) error {
	class, err := models.ClassByID(s, classId)
	if err != nil {
		return err
	}
	err = class.Delete(s)
	return err
}

func (s *postgresService) GetClassMembers(classId uuid.UUID) ([]*models.Member, error) {
	return models.MembersByClassID(s, classId)
}

func (s *postgresService) GetMember(id uuid.UUID) (*models.Member, error) {
	return models.MemberByID(s, id)
}

func (s *postgresService) GetMemberByUser(userId uuid.UUID, classId uuid.UUID) (*models.Member, error) {
	return models.MemberByUserIDClassID(s, userId, classId)
}

func (s *postgresService) CreateMember(member *models.Member) error {
	return member.Save(s)
}

func (s *postgresService) DeleteMember(id uuid.UUID) error {
	member, err := models.MemberByID(s, id)
	if err != nil {
		return err
	}
	err = member.Delete(s)
	return err
}

func (s *postgresService) UpdateMember(member *models.Member) error {
	return member.Update(s)
}

func NewPostgres(db *sql.DB) Service {
	return &postgresService{db}
}
