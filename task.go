package main

type Task struct {
	ID      int64
	Payload *OpenPositionParams
}

func NewTask(id int64) *Task {
	return &Task{
		ID:      id,
		Payload: NewOpenPositionParams(0, 0.02, 0),
	}
}
