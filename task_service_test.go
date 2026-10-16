package gosis1

import "testing"

func TestTaskLifecycle(t *testing.T) {
	service := NewTaskService(&DefaultURLParser{})

	task, err := service.CreateTask(URLInput{"https://test.test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Status != StatusPending {
		t.Errorf("expected status PENDING, got %v", task.Status)
	}

	if _, err = service.GetTask("asdfasdfadf"); err == nil {
		t.Error("expected error for non-existent task, got nil")
	}

	processed, err := service.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("failed to process task: %v", err)
	}
	if processed.Status != StatusCompleted {
		t.Errorf("expected status COMPLETED, got %v", processed.Status)
	}

	_, err = service.ProcessTask(task.ID)
	if err == nil {
		t.Error("expected error when re-processing completed task, got nil")
	}
}
