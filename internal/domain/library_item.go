package domain

import (
	"errors"
	"time"
)

var ErrLibraryItemNotFound = errors.New("library item not found")

type Status string

const (
	StatusPlanning   Status = "planning"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusDropped    Status = "dropped"
	StatusOnHold     Status = "on_hold"
)

type LibraryItem struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Media      Media     `json:"media"`
	Status     Status    `json:"status"`
	UserRating *int      `json:"user_rating,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	StartedAt  *string   `json:"started_at,omitempty"`
	FinishedAt *string   `json:"finished_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
