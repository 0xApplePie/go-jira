package models

import (
	"time"
	"errors"
)

type Status string

const (
	StatusTodo     Status = "TODO"
	StatusProgress Status = "PROGRESS"
	StatusDone     Status = "DONE"
)

func ParseStatus(s string) (Status, error) {
	switch s {
	case string(StatusTodo):
		return StatusTodo, nil
	case string(StatusProgress):
		return StatusProgress, nil
	case string(StatusDone):
		return StatusDone, nil
	default:
		return "", errors.New("invalid status")
	}
}

type Ticket struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Assignee    *string   `json:"assignee,omitempty"`
} 