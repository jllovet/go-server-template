package todo_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/jllovet/go-server-template/internal/todo"
)

var errRepository = errors.New("repository error")

// mockRepository is a mock implementation of the todo.Repository for testing.
type mockRepository struct {
	mu    sync.RWMutex
	todos map[string]todo.Todo

	// Control fields to simulate errors
	saveErr     error
	findByIDErr error
	findAllErr  error
	deleteErr   error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		todos: make(map[string]todo.Todo),
	}
}

func (m *mockRepository) Save(_ context.Context, t todo.Todo) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.saveErr != nil {
		return m.saveErr
	}
	m.todos[t.ID] = t
	return nil
}

func (m *mockRepository) FindByID(_ context.Context, id string) (todo.Todo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.findByIDErr != nil {
		return todo.Todo{}, m.findByIDErr
	}
	t, ok := m.todos[id]
	if !ok {
		return todo.Todo{}, errors.New("not found")
	}
	return t, nil
}

func (m *mockRepository) FindAll(_ context.Context) ([]todo.Todo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	all := make([]todo.Todo, 0, len(m.todos))
	for _, t := range m.todos {
		all = append(all, t)
	}
	return all, nil
}

func (m *mockRepository) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.todos[id]; !ok {
		return errors.New("not found")
	}
	delete(m.todos, id)
	return nil
}

func TestService(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		repo := newMockRepository()
		service := todo.NewService(repo)

		// Success case
		created, err := service.Create(ctx, "Test Todo")
		if err != nil {
			t.Fatalf("Create() error = %v, want nil", err)
		}
		if created.Title != "Test Todo" {
			t.Errorf("Create() got title = %q, want %q", created.Title, "Test Todo")
		}
		if created.ID == "" {
			t.Error("Create() got empty ID, want non-empty")
		}

		// Validation error
		_, err = service.Create(ctx, "")
		if err == nil {
			t.Fatal("Create() with empty title expected error, got nil")
		}

		// Repository error
		repo.saveErr = errRepository
		_, err = service.Create(ctx, "Another Todo")
		if !errors.Is(err, errRepository) {
			t.Fatalf("Create() expected repository error, got %v", err)
		}
	})

	t.Run("List", func(t *testing.T) {
		repo := newMockRepository()
		service := todo.NewService(repo)

		// Seed data
		_, _ = service.Create(ctx, "First")
		_, _ = service.Create(ctx, "Second")

		// Success case
		todos, err := service.List(ctx)
		if err != nil {
			t.Fatalf("List() error = %v, want nil", err)
		}
		if len(todos) != 2 {
			t.Fatalf("List() got %d todos, want 2", len(todos))
		}

		// Repository error
		repo.findAllErr = errRepository
		_, err = service.List(ctx)
		if !errors.Is(err, errRepository) {
			t.Fatalf("List() expected repository error, got %v", err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		repo := newMockRepository()
		service := todo.NewService(repo)
		created, _ := service.Create(ctx, "My Todo")

		// Success case
		found, err := service.Get(ctx, created.ID)
		if err != nil {
			t.Fatalf("Get() error = %v, want nil", err)
		}
		if !reflect.DeepEqual(found, created) {
			t.Errorf("Get() got = %v, want %v", found, created)
		}

		// Not found
		_, err = service.Get(ctx, "non-existent-id")
		if err == nil {
			t.Fatal("Get() with non-existent ID expected error, got nil")
		}
	})

	t.Run("Update", func(t *testing.T) {
		repo := newMockRepository()
		service := todo.NewService(repo)
		created, _ := service.Create(ctx, "Original Title")

		// Success case: update title
		newTitle := "Updated Title"
		updated, err := service.Update(ctx, created.ID, newTitle)
		if err != nil {
			t.Fatalf("Update() error = %v, want nil", err)
		}
		if updated.Title != newTitle {
			t.Errorf("Update() got = %q, want %q", updated.Title, newTitle)
		}
		if updated.Completed != created.Completed {
			t.Error("Update() should not change completed status")
		}

		// Not found
		_, err = service.Update(ctx, "non-existent-id", newTitle)
		if err == nil {
			t.Fatal("Update() with non-existent ID expected error, got nil")
		}

		// Validation error
		_, err = service.Update(ctx, created.ID, "")
		if err == nil {
			t.Fatal("Update() with empty title expected error, got nil")
		}

		// Repository error
		repo.saveErr = errRepository
		anotherTitle := "Final Title"
		_, err = service.Update(ctx, created.ID, anotherTitle)
		if !errors.Is(err, errRepository) {
			t.Fatalf("Update() expected repository error, got %v", err)
		}
	})

	t.Run("SetCompleted", func(t *testing.T) {
		repo := newMockRepository()
		service := todo.NewService(repo)
		created, _ := service.Create(ctx, "Original Title")

		// Success case: update completed
		newCompleted := true
		updated, err := service.SetCompleted(ctx, created.ID, newCompleted)
		if err != nil {
			t.Fatalf("SetCompleted() error = %v, want nil", err)
		}
		if updated.Completed != newCompleted {
			t.Errorf("SetCompleted() got = %v, want %v", updated.Completed, newCompleted)
		}
		if updated.Title != created.Title {
			t.Errorf("SetCompleted() should not change title, got %q", updated.Title)
		}

		// Not found
		_, err = service.SetCompleted(ctx, "non-existent-id", true)
		if err == nil {
			t.Fatal("SetCompleted() with non-existent ID expected error, got nil")
		}

		// Repository error
		repo.saveErr = errRepository
		_, err = service.SetCompleted(ctx, created.ID, false)
		if !errors.Is(err, errRepository) {
			t.Fatalf("SetCompleted() expected repository error, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		repo := newMockRepository()
		service := todo.NewService(repo)
		created, _ := service.Create(ctx, "To Be Deleted")

		// Success case
		err := service.Delete(ctx, created.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v, want nil", err)
		}

		// Verify it's gone
		_, err = service.Get(ctx, created.ID)
		if err == nil {
			t.Fatal("Get() after Delete() expected error, got nil")
		}

		// Not found
		err = service.Delete(ctx, "non-existent-id")
		if err == nil {
			t.Fatal("Delete() with non-existent ID expected error, got nil")
		}

		// Repository error
		created2, _ := service.Create(ctx, "Another one")
		repo.deleteErr = errRepository
		err = service.Delete(ctx, created2.ID)
		if !errors.Is(err, errRepository) {
			t.Fatalf("Delete() expected repository error, got %v", err)
		}
	})
}
