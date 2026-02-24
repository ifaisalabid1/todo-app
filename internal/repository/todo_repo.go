package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ifaisalabid1/todo-app/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *domain.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	GetAll(ctx context.Context) ([]*domain.Todo, error)
	Update(ctx context.Context, todo *domain.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type todoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoRepository(pool *pgxpool.Pool) TodoRepository {
	return &todoRepository{pool: pool}
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `
			INSERT INTO todos (id, title, description, status, due_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

	now := time.Now().UTC()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	args := []any{
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.DueDate,
		todo.CreatedAt,
		todo.UpdatedAt,
	}

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrTodoAlreadyExists
			}
		}

		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	query := `
			SELECT id, title, description, status, due_date, created_at, updated_at
			FROM todos WHERE id = $1
	`

	var todo domain.Todo

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.DueDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTodoNotFound
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

func (r *todoRepository) GetAll(ctx context.Context) ([]*domain.Todo, error) {
	query := `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	defer rows.Close()

	var todos []*domain.Todo

	for rows.Next() {
		var todo domain.Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.DueDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}

		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
		UPDATE todos
		SET title = $2, description = $3, status = $4, due_date = $5, updated_at = $6
		WHERE id = $1
	`

	todo.UpdatedAt = time.Now().UTC()

	result, err := r.pool.Exec(ctx, query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.DueDate,
		todo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTodoNotFound
	}

	return nil

}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM todos WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTodoNotFound
	}

	return nil
}

func (r *todoRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1)"

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check todo existence: %w", err)
	}

	return exists, nil
}
