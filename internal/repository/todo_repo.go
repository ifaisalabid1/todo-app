package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ifaisalabid1/todo-app/internal/domain"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *domain.CreateTodoInput) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	GetAll(ctx context.Context) ([]*domain.Todo, error)
	Update(ctx context.Context, todo *domain.UpdateTodoInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
