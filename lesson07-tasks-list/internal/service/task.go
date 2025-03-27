package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"tasks/internal/repository"
)

type Service struct {
	DB repository.DB
}

func NewService(ctx context.Context, db repository.DB) (*Service, error) {
	s := &Service{
		DB: db,
	}

	err := s.DB.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) AddTask(ctx context.Context, name string) error {
	err := s.DB.AddTask(ctx, name)
	if err != nil {
		return fmt.Errorf("adding task: %w", err)
	}

	fmt.Printf("Added task: %s\n", name)

	return nil
}

func (s *Service) ListOutStandingTasks(ctx context.Context) error {
	results, err := s.DB.ListOutstandingTasks(ctx)
	if err != nil {
		return fmt.Errorf("listing out standing tasks: %w", err)
	}

	fmt.Println("Listing outstanding tasks:")

	// Define column widths
	idWidth := 5
	nameWidth := 20
	statusWidth := 6

	// Pad with spaces on the right instead of left to align the fields
	fmt.Printf("%-*s %-*s %-*s\n", idWidth, "ID", nameWidth, "Task", statusWidth, "Completed")
	fmt.Println(strings.Repeat("-", idWidth+nameWidth+statusWidth+5))
	for _, task := range results {
		fmt.Printf("%-*d %-*s %-*t\n", idWidth, task.ID, nameWidth, task.Name, statusWidth, task.Completed)
	}

	return nil
}

func (s *Service) CompleteTask(ctx context.Context, taskID string) error {
	task, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("completing task: converting taskID %s to an int", taskID)
	}

	err = s.DB.CompleteTask(ctx, task)
	if err != nil {
		return fmt.Errorf("completing task: %w", err)
	}

	fmt.Printf("Completed task: %d", task)

	return nil
}
