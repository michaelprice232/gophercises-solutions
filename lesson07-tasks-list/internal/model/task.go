package model

import "time"

type Task struct {
	ID          int
	Name        string
	Completed   bool
	CompletedAt time.Time
}

type Tasks []Task
