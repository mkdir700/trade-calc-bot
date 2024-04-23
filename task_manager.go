package main

import (
	"sync"
)

var taskManager *TaskManager

type TaskManager struct {
	tasks sync.Map
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
