package model

type Task struct {
	ID        int
	Name      string
	Completed bool
}

type Tasks []Task
