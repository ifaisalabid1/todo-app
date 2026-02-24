package domain

import (
	"time"
)

type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted  TodoStatus = "completed"
)

type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TodoStatus `json:"status"`
	DueDate     time.Time  `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTodoInput struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"required,min=1,max=500"`
	Status      TodoStatus `json:"status"`
	DueDate     time.Time  `json:"due_date,omitzero"`
}

type UpdateTodoInput struct {
	Title       *string     `json:"title" validate:"required,min=1,max=255"`
	Description *string     `json:"description" validate:"required,min=1,max=500"`
	Status      *TodoStatus `json:"status"`
	DueDate     *time.Time  `json:"due_date,omitzero"`
}
