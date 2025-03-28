package repository

import (
	"context"
	"tasks/internal/model"
)

type DB interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	AddTask(ctx context.Context, name string) error
	CompleteTask(ctx context.Context, id int) error
	ListOutstandingTasks(ctx context.Context) (model.Tasks, error)
	ListCompletedTasks(ctx context.Context) (model.Tasks, error)
}
