package server

import (
	"net/http"
)

func (s *Server) handleCreateTodo() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := s.decode(w, r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := s.service.Create(r.Context(), req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.encode(w, http.StatusCreated, t)
	}
}

func (s *Server) handleListTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := s.service.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.encode(w, http.StatusOK, todos)
	}
}

func (s *Server) handleGetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		t, err := s.service.Get(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		s.encode(w, http.StatusOK, t)
	}
}

func (s *Server) handleUpdateTodoTitle() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req request
		if err := s.decode(w, r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := s.service.Update(r.Context(), id, req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.encode(w, http.StatusOK, t)
	}
}

func (s *Server) handleMarkTodoComplete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		t, err := s.service.SetCompleted(r.Context(), id, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.encode(w, http.StatusOK, t)
	}
}

func (s *Server) handleMarkTodoIncomplete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		t, err := s.service.SetCompleted(r.Context(), id, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.encode(w, http.StatusOK, t)
	}
}

func (s *Server) handleDeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := s.service.Delete(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
