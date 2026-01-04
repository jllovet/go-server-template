package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/jllovet/go-server-template/internal/todo"
)

// Repository is an in-memory implementation of todo.Repository.
type Repository struct {
	// mu protects the todos map from concurrent access.
	mu    sync.RWMutex
	todos map[string]todo.Todo
}

// New creates a new in-memory repository.
func New() *Repository {
	return &Repository{
		todos: make(map[string]todo.Todo),
	}
}

// Save stores the todo item.
func (r *Repository) Save(ctx context.Context, t todo.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.todos[t.ID] = t
	return nil
}

// FindByID retrieves a todo by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.todos[id]
	if !ok {
		return todo.Todo{}, fmt.Errorf("todo not found")
	}
	return t, nil
}

// FindAll retrieves all todos.
func (r *Repository) FindAll(ctx context.Context) ([]todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	todos := make([]todo.Todo, 0, len(r.todos))
	for _, t := range r.todos {
		todos = append(todos, t)
	}
	return todos, nil
}

// Delete removes a todo by its ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.todos[id]; !ok {
		return fmt.Errorf("todo not found")
	}
	delete(r.todos, id)
	return nil
}
