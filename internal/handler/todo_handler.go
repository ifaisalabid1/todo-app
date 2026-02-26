package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ifaisalabid1/todo-app/internal/domain"
	"github.com/ifaisalabid1/todo-app/internal/service"
	"github.com/ifaisalabid1/todo-app/pkg/response"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) RegisterRoutes(router chi.Router) {
	router.Route("/todos", func(r chi.Router) {
		r.Get("/", h.GetAllTodos)
		r.Post("/", h.CreateTodo)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetTodoByID)
			r.Put("/", h.UpdateTodo)
			r.Patch("/", h.PartialUpdateTodo)
			r.Delete("/", h.DeleteTodo)
		})
	})
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	todo, err := h.service.CreateTodo(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrTodoAlreadyExists) {
			response.Error(w, http.StatusConflict, "todo already exists", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create todo", err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "todo created successfully", todo)
}

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid todo id", err.Error())
		return
	}

	todo, err := h.service.GetTodoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get todo", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "todo received successfully", todo)
}

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get todos", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "todos received successfully", todos)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid todo id", err.Error())
		return
	}

	var input domain.UpdateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	todo, err := h.service.UpdateTodo(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update todo", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "todo updated successfully", todo)
}

func (h *TodoHandler) PartialUpdateTodo(w http.ResponseWriter, r *http.Request) {
	h.UpdateTodo(w, r)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid todo id", err.Error())
		return
	}

	if err := h.service.DeleteTodo(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete todo", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "todo deleted successfully", nil)
}
