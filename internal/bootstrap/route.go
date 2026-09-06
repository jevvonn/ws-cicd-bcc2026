package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/handler"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/repository"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/usecase"
	"github.com/jevvonn/ws-cicd-bcc2026/pkg/helper"
	"gorm.io/gorm"
)

func registerRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello, welcome to the API of ws-cicd-bcc2026!"})
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "hello", fiber.Map{
			"app":  helper.Env("APP_NAME", "ws-cicd-bcc2026"),
			"env":  helper.Env("APP_ENV", "development"),
			"text": "Hello, World!",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(c.Context()) != nil {
			return helper.Error(c, fiber.StatusServiceUnavailable, "database unreachable")
		}

		return helper.Success(c, fiber.StatusOK, "healthy", fiber.Map{
			"status":   "ok",
			"database": "ok",
		})
	})

	todoRepository := repository.NewTodoRepository(db)
	todoUsecase := usecase.NewTodoUsecase(todoRepository)

	api := app.Group("/api")
	handler.NewTodoHandler(api, todoUsecase)
}
