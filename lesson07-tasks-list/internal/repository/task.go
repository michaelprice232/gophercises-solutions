package repository

import (
	"gophercises-ex7/internal/model"
)

type DB interface {
	Connect() error
	AddTask(name string) error
	CompleteTask(id int) error
	ListOutstandingTasks() (model.Tasks, error)
}
