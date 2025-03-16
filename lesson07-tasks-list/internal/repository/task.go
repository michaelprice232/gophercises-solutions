package repository

import (
	"tasks/internal/model"
)

type DB interface {
	Connect() error
	Close() error
	AddTask(name string) error
	CompleteTask(id int) error
	ListOutstandingTasks() (model.Tasks, error)
}
