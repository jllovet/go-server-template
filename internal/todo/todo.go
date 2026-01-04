package todo

import "context"

// Todo represents a task in the system.
type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Repository defines the interface for storing and retrieving Todos.
// In Hexagonal Architecture, this is a "Driven Port".
type Repository interface {
	Save(ctx context.Context, t Todo) error
	FindByID(ctx context.Context, id string) (Todo, error)
	FindAll(ctx context.Context) ([]Todo, error)
	Delete(ctx context.Context, id string) error
}

// Service defines the interface for the business logic.
// In Hexagonal Architecture, this is a "Driving Port" used by the HTTP handler.
type Service interface {
	Create(ctx context.Context, title string) (Todo, error)
	Get(ctx context.Context, id string) (Todo, error)
	List(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, id string, title string) (Todo, error)
	SetCompleted(ctx context.Context, id string, completed bool) (Todo, error)
	Delete(ctx context.Context, id string) error
}
