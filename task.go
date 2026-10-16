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
	ErrTas