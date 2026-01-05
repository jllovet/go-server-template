package memory_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jllovet/go-server-template/internal/todo"
	"github.com/jllovet/go-server-template/internal/todo/memory"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("Save and FindByID", func(t *testing.T) {
		repo := memory.New()
		item := todo.Todo{
			ID:    "1",
			Title: "Test Item",
		}

		if err := repo.Save(ctx, item); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		found, err := repo.FindByID(ctx, "1")
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found.Title != item.Title {
			t.Errorf("got title %q, want %q", found.Title, item.Title)
		}
	})

	t.Run("FindByID Not Found", func(t *testing.T) {
		repo := memory.New()
		_, err := repo.FindByID(ctx, "non-existent")
		if err == nil {
			t.Error("FindByID() expected error for non-existent item, got nil")
		}
	})

	t.Run("Update", func(t *testing.T) {
		repo := memory.New()
		item := todo.Todo{ID: "1", Title: "Original"}
		_ = repo.Save(ctx, item)

		updated := todo.Todo{ID: "1", Title: "Updated"}
		if err := repo.Save(ctx, updated); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		found, _ := repo.FindByID(ctx, "1")
		if found.Title != "Updated" {
			t.Errorf("got title %q, want %q", found.Title, "Updated")
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		repo := memory.New()
		item1 := todo.Todo{ID: "1", Title: "Item 1"}
		item2 := todo.Todo{ID: "2", Title: "Item 2"}
		_ = repo.Save(ctx, item1)
		_ = repo.Save(ctx, item2)

		all, err := repo.FindAll(ctx)
		if err != nil {
			t.Fatalf("FindAll() error = %v", err)
		}

		if len(all) != 2 {
			t.Errorf("got %d items, want 2", len(all))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		repo := memory.New()
		item := todo.Todo{ID: "3", Title: "To Delete"}
		_ = repo.Save(ctx, item)

		if err := repo.Delete(ctx, "3"); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		_, err := repo.FindByID(ctx, "3")
		if err == nil {
			t.Error("FindByID() expected error after delete, got nil")
		}
	})

	t.Run("Delete Not Found", func(t *testing.T) {
		repo := memory.New()
		if err := repo.Delete(ctx, "non-existent"); err == nil {
			t.Error("Delete() expected error for non-existent item, got nil")
		}
	})

	t.Run("Concurrent Access", func(t *testing.T) {
		// This test verifies that the repository is thread-safe.
		// Go maps are not safe for concurrent use, so this test would panic
		// without the sync.RWMutex in the repository implementation.
		repo := memory.New()
		var wg sync.WaitGroup
		count := 100

		// Concurrent writes
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(i int) {
				defer wg.Done()
				id := fmt.Sprintf("%d", i)
				_ = repo.Save(ctx, todo.Todo{ID: id, Title: "Concurrent"})
			}(i)
		}

		// Concurrent reads
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(i int) {
				defer wg.Done()
				id := fmt.Sprintf("%d", i)
				_, _ = repo.FindByID(ctx, id)
			}(i)
		}

		wg.Wait()

		all, err := repo.FindAll(ctx)
		if err != nil {
			t.Fatalf("FindAll() error = %v", err)
		}
		if len(all) != count {
			t.Errorf("got %d items, want %d", len(all), count)
		}
	})
}

func BenchmarkRepository_Save(b *testing.B) {
	ctx := context.Background()
	repo := memory.New()
	item := todo.Todo{ID: "1", Title: "Benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.Save(ctx, item)
	}
}

func BenchmarkRepository_FindByID(b *testing.B) {
	ctx := context.Background()
	repo := memory.New()
	item := todo.Todo{ID: "1", Title: "Benchmark"}
	_ = repo.Save(ctx, item)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.FindByID(ctx, "1")
	}
}
