package service

import (
	"gophercises-ex7/internal/repository"
)

type Service struct {
	db repository.DB
}

func NewService(db repository.DB) (*Service, error) {
	return &Service{
		db: db,
	}, nil
}

func (s *Service) AddTask(name string) error {
	// Perform other business logic here

	return s.db.AddTask(name)
}
