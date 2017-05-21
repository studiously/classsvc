package service

import "github.com/google/uuid"

type Service interface {
}

type Class struct {
	Id      uuid.UUID
	Members []uuid.UUID
}

type ClassMember struct {
	UserId uuid.UUID
	
}
