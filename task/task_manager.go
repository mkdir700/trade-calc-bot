package task

import (
	"sync"
)

var taskManager *TaskManager

type TaskManager struct {
	tasks sync.Map
}

// InputStep represents a step in the input process.
type InputStep int

const (
	Start InputStep = iota
	InputCapital
	InputCapitalLossRadio
	InputLossRadio
	End
)

type Task struct {
	Payload *OpenPosition
	ID      int64
	Step    InputStep
}

func NewTask(id int64) *Task {
	return &Task{
		ID:      id,
		Step:    Start,
		Payload: &OpenPosition{},
	}
}

func (t *Task) NextStep() {
	t.Step++
}

func (tm *TaskManager) AddTask(task *Task) {
	tm.tasks.Store(task.ID, task)
}

func (tm *TaskManager) GetTask(id int64) *Task {
	if task, ok := tm.tasks.Load(id); ok {
		return task.(*Task)
	}
	return nil
}

func (tm *TaskManager) HasTask(id int64) bool {
	_, ok := tm.tasks.Load(id)
	return ok
}

func (tm *TaskManager) RemoveTask(id int64) {
	tm.tasks.Delete(id)
}

func GetTaskManager() *TaskManager {
	return taskManager
}

func init() {
	taskManager = &TaskManager{}
}
