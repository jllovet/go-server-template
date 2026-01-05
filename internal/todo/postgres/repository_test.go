package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jllovet/go-server-template/internal/todo"
	"github.com/jllovet/go-server-template/internal/todo/postgres"
)

func TestRepository(t *testing.T) {
	// Skip if TEST_DATABASE_URL is not set.
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("Skipping postgres repository tests: TEST_DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	// Ensure table exists for tests
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE
		);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Helper to clean DB between tests
	cleanDB := func() {
		_, err := db.Exec("TRUNCATE TABLE todos")
		if err != nil {
			t.Fatalf("failed to truncate table: %v", err)
		}
	}

	ctx := context.Background()

	t.Run("Save and FindByID", func(t *testing.T) {
		cleanDB()
		repo := postgres.New(db)
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
		if found.Completed != false {
			t.Errorf("got completed %v, want false", found.Completed)
		}
	})

	t.Run("FindByID Not Found", func(t *testing.T) {
		cleanDB()
		repo := postgres.New(db)
		_, err := repo.FindByID(ctx, "non-existent")
		if err == nil {
			t.Error("FindByID() expected error for non-existent item, got nil")
		}
	})

	t.Run("Update (Upsert)", func(t *testing.T) {
		cleanDB()
		repo := postgres.New(db)
		item := todo.Todo{ID: "1", Title: "Original", Completed: false}
		if err := repo.Save(ctx, item); err != nil {
			t.Fatalf("Save() initial error = %v", err)
		}

		updated := todo.Todo{ID: "1", Title: "Updated", Completed: true}
		if err := repo.Save(ctx, updated); err != nil {
			t.Fatalf("Save() update error = %v", err)
		}

		found, err := repo.FindByID(ctx, "1")
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found.Title != "Updated" {
			t.Errorf("got title %q, want %q", found.Title, "Updated")
		}
		if !found.Completed {
			t.Error("got completed false, want true")
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		cleanDB()
		repo := postgres.New(db)
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
		cleanDB()
		repo := postgres.New(db)
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
		cleanDB()
		repo := postgres.New(db)
		if err := repo.Delete(ctx, "non-existent"); err == nil {
			t.Error("Delete() expected error for non-existent item, got nil")
		}
	})
}
