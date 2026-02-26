package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ifaisalabid1/todo-app/internal/domain"
	"github.com/ifaisalabid1/todo-app/internal/repository"
)

type TodoService struct {
	repo      repository.TodoRepository
	validator *validator.Validate
	logger    *slog.Logger
}

func NewTodoService(repo repository.TodoRepository, logger *slog.Logger) *TodoService {
	return &TodoService{
		repo:      repo,
		validator: validator.New(),
		logger:    logger,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, input domain.CreateTodoInput) (*domain.Todo, error) {
	s.logger.InfoContext(ctx, "creating todo", "title", input.Title)

	if err := s.validator.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	todo := &domain.Todo{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		DueDate:     input.DueDate,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		s.logger.ErrorContext(ctx, "failed to create todo", "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "todo created successfully", "todo_id", todo.ID)
	return todo, nil
}

func (s *TodoService) GetTodoByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	s.logger.InfoContext(ctx, "getting todo", "todo_id", id)

	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get todo", "error", err)
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]*domain.Todo, error) {
	s.logger.InfoContext(ctx, "getting all todos")

	todos, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get todos", "error", err)
		return nil, err
	}

	return todos, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, id uuid.UUID, input domain.UpdateTodoInput) (*domain.Todo, error) {
	s.logger.InfoContext(ctx, "updating todo", "todo_id", id)

	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get todo", "error", err)
		return nil, err
	}

	if input.Title != nil {
		todo.Title = *input.Title
	}

	if input.Description != nil {
		todo.Description = *input.Description
	}

	if input.Status != nil {
		todo.Status = *input.Status
	}

	if input.DueDate != nil {
		todo.DueDate = *input.DueDate
	}

	if err := s.validator.Struct(todo); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Update(ctx, todo); err != nil {
		s.logger.ErrorContext(ctx, "failed to update todo", "error", err)
		return nil, err
	}

	s.logger.InfoContext(ctx, "todo updated successfully", "todo_id", id)
	return todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	s.logger.InfoContext(ctx, "deleting todo", "todo_id", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			s.logger.ErrorContext(ctx, "todo not found for deletion", "todo_id", id)
			return err
		}

		s.logger.ErrorContext(ctx, "failed to delete toto", "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "todo deleted successfully", "todo_id", id)
	return nil

}
