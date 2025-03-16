package service

import (
	"log/slog"
	"tasks/internal/repository"
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
	slog.Info("Adding task.", "task", name)

	err := s.db.Connect()
	if err != nil {
		return err
	}

	defer func() {
		if err = s.db.Close(); err != nil {
			slog.Error("problem closing database connection", "error", err)
		}
	}()

	return s.db.AddTask(name)
}
