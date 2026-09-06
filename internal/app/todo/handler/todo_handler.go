package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/entity"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/usecase"
	"github.com/jevvonn/ws-cicd-bcc2026/pkg/helper"
)

type TodoHandler struct {
	usecase usecase.TodoUsecase
}

func NewTodoHandler(router fiber.Router, todoUsecase usecase.TodoUsecase) {
	h := &TodoHandler{usecase: todoUsecase}

	todos := router.Group("/todos")
	todos.Post("/", h.Create)
	todos.Get("/", h.GetAll)
	todos.Get("/:id", h.GetByID)
	todos.Put("/:id", h.Update)
	todos.Delete("/:id", h.Delete)
}

func (h *TodoHandler) Create(c *fiber.Ctx) error {
	var req entity.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	todo, err := h.usecase.Create(req)
	if err != nil {
		return mapError(c, err)
	}

	return helper.Success(c, fiber.StatusCreated, "todo created", todo)
}

func (h *TodoHandler) GetAll(c *fiber.Ctx) error {
	todos, err := h.usecase.GetAll()
	if err != nil {
		return mapError(c, err)
	}

	return helper.Success(c, fiber.StatusOK, "todos retrieved", todos)
}

func (h *TodoHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	todo, err := h.usecase.GetByID(id)
	if err != nil {
		return mapError(c, err)
	}

	return helper.Success(c, fiber.StatusOK, "todo retrieved", todo)
}

func (h *TodoHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	var req entity.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	todo, err := h.usecase.Update(id, req)
	if err != nil {
		return mapError(c, err)
	}

	return helper.Success(c, fiber.StatusOK, "todo updated", todo)
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.usecase.Delete(id); err != nil {
		return mapError(c, err)
	}

	return helper.Success(c, fiber.StatusOK, "todo deleted", nil)
}

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func mapError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrTodoNotFound):
		return helper.Error(c, fiber.StatusNotFound, err.Error())
	case errors.Is(err, usecase.ErrTitleRequired):
		return helper.Error(c, fiber.StatusBadRequest, err.Error())
	default:
		return helper.Error(c, fiber.StatusInternalServerError, "internal server error")
	}
}
