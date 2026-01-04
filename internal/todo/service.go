package todo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// service implements the Service interface.
// It holds a reference to the Repository Port.
type service struct {
	repo Repository
}

// NewService creates a new Todo service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create applies business logic to create a new Todo.
func (s *service) Create(ctx context.Context, title string) (Todo, error) {
	if title == "" {
		return Todo{}, fmt.Errorf("title cannot be empty")
	}

	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return Todo{}, fmt.Errorf("failed to generate id: %w", err)
	}

	t := Todo{
		ID:    hex.EncodeToString(id),
		Title: title,
	}

	if err := s.repo.Save(ctx, t); err != nil {
		return Todo{}, fmt.Errorf("failed to save todo: %w", err)
	}

	return t, nil
}

func (s *service) List(ctx context.Context) ([]Todo, error) {
	todos, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	return todos, nil
}

func (s *service) Get(ctx context.Context, id string) (Todo, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Todo{}, fmt.Errorf("failed to get todo %q: %w", id, err)
	}
	return t, nil
}

func (s *service) Update(ctx context.Context, id string, title string) (Todo, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Todo{}, fmt.Errorf("failed to find todo for update: %w", err)
	}

	if title == "" {
		return Todo{}, fmt.Errorf("title cannot be empty")
	}
	t.Title = title

	if err := s.repo.Save(ctx, t); err != nil {
		return Todo{}, fmt.Errorf("failed to save updated todo: %w", err)
	}

	return t, nil
}

func (s *service) SetCompleted(ctx context.Context, id string, completed bool) (Todo, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Todo{}, fmt.Errorf("failed to find todo for update: %w", err)
	}

	t.Completed = completed

	if err := s.repo.Save(ctx, t); err != nil {
		return Todo{}, fmt.Errorf("failed to save updated todo: %w", err)
	}

	return t, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete todo %q: %w", id, err)
	}
	return nil
}
