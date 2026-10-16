package gosis1

import (
	"fmt"
	"time"
)

func (s *TaskService) CreateTask(input URLInput) (*Task, error) {
	if _, err := s.parser.Parse(input); err != nil {
		return nil, err
	}

	id := fmt.Sprintf("task-%d", s.nextId)
	s.nextId++

	now := time.Now()
	task := &Task{
		ID:        id,
		Status:    StatusPending,
		Input:     input,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.tasks[id] = task
	return task, nil
}

func (s *TaskService) ProcessTask(id string) (*Task, error) {
	task, existent := s.tasks[id]
	if !existent {
		return nil, ErrTaskNotFound
	}

	if task.Status == StatusCompleted || task.Status == StatusFailed {
		return nil, ErrTaskAlreadyHandled
	}

	task.Status = StatusProcessing
	task.UpdatedAt = time.Now()

	res, err := s.parser.Parse(task.Input)
	if err != nil {
		task.Status = StatusFailed
		task.Errormsg = err.Error()
		task.UpdatedAt = time.Now()
		return nil, err
	}

	task.Result = res
	task.Status = StatusCompleted
	task.UpdatedAt = time.Now()

	return task, nil
}

func (s *TaskService) GetTask(id string) (*Task, error) {
	task, existent := s.tasks[id]
	if !existent {
		return nil, ErrTaskNotFound
	}
	return task, nil
}
