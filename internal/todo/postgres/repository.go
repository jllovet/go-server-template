package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jllovet/go-server-template/internal/todo"
)

// Repository implements todo.Repository using PostgreSQL.
type Repository struct {
	db *sql.DB
}

// New creates a new Postgres repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save creates or updates a todo.
func (r *Repository) Save(ctx context.Context, t todo.Todo) error {
	query := `
		INSERT INTO todos (id, title, completed)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET title = EXCLUDED.title, completed = EXCLUDED.completed
	`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.Title, t.Completed)
	if err != nil {
		return fmt.Errorf("postgres save: %w", err)
	}
	return nil
}

// FindByID retrieves a todo by ID.
func (r *Repository) FindByID(ctx context.Context, id string) (todo.Todo, error) {
	query := `SELECT id, title, completed FROM todos WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var t todo.Todo
	if err := row.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo.Todo{}, fmt.Errorf("todo not found")
		}
		return todo.Todo{}, fmt.Errorf("postgres find by id: %w", err)
	}
	return t, nil
}

// FindAll retrieves all todos.
func (r *Repository) FindAll(ctx context.Context) ([]todo.Todo, error) {
	query := `SELECT id, title, completed FROM todos`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres find all: %w", err)
	}
	defer rows.Close()

	var todos []todo.Todo
	for rows.Next() {
		var t todo.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			return nil, fmt.Errorf("postgres scan: %w", err)
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// Delete removes a todo by ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM todos WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres delete: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}
