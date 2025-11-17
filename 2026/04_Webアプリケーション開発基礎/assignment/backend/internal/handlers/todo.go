package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/store"

	"github.com/labstack/echo/v4"
)

// TodoHandler はTODO関連のHTTPハンドラーを提供する構造体
type TodoHandler struct {
	store *store.TodoStore
}

// NewTodoHandler は新しいTodoHandlerインスタンスを作成
func NewTodoHandler(store *store.TodoStore) *TodoHandler {
	return &TodoHandler{store: store}
}

// ListTodos は全てのTODOを取得するハンドラー
func (h *TodoHandler) ListTodos(c echo.Context) error {
	return c.JSON(http.StatusOK, h.store.List())
}

// CreateTodo は新しいTODOを作成するハンドラー
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var req models.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	todo, err := h.store.Create(req.Title)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, todo)
}

// UpdateTodo は既存のTODOを更新するハンドラー
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var req models.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Title == nil && req.Completed == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	todo, err := h.store.Update(id, req)
	if err != nil {
		if errors.Is(err, echo.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

// DeleteTodo は指定されたTODOを削除するハンドラー
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, echo.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes はTODOハンドラーのルートをEchoインスタンスに登録
func (h *TodoHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/todos", h.ListTodos)
	e.POST("/todos", h.CreateTodo)
	e.PATCH("/todos/:id", h.UpdateTodo)
	e.DELETE("/todos/:id", h.DeleteTodo)
}
