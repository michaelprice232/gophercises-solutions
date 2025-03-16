package service

import (
	"fmt"
	"strings"
	"tasks/internal/repository"
)

type Service struct {
	DB repository.DB
}

func NewService(db repository.DB) (*Service, error) {
	s := &Service{
		DB: db,
	}

	err := s.DB.Connect()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) AddTask(name string) error {
	fmt.Printf("Added task: %s\n", name)

	return s.DB.AddTask(name)
}

func (s *Service) ListOutStandingTasks() error {
	fmt.Println("Listing out standing tasks:")

	results, err := s.DB.ListOutstandingTasks()
	if err != nil {
		return fmt.Errorf("listing out standing tasks: %w", err)
	}

	// Define column widths
	idWidth := 5
	nameWidth := 20
	statusWidth := 6

	fmt.Printf("%-*s %-*s %-*s\n", idWidth, "ID", nameWidth, "Task", statusWidth, "Completed")
	fmt.Println(strings.Repeat("-", idWidth+nameWidth+statusWidth+5))
	for _, task := range results {
		fmt.Printf("%-*d %-*s %-*t\n", idWidth, task.ID, nameWidth, task.Name, statusWidth, task.Completed)
	}

	return nil
}
