package gosis1

import (
	"errors"
	"time"
)

type TaskStatus string

const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyHandled = errors.New("task already handled")
	ErrInvaidInput        = errors.New("invalid input")
	ErrInvalidTransition  = errors.New("invalid transition")
)

type URLInput struct {
	RawURL string
}

type ParsedPathResult struct {
	Scheme       string
	Host         string
	FullPath     string
	Segments     []string
	SegmentCount int
}

type Task struct {
	ID        string
	Status    TaskStatus
	Input     URLInput
	Result    *ParsedPathResult
	Errormsg  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
