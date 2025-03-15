package model

type Task struct {
	id        int
	name      string
	completed bool
}

type Tasks []Task
