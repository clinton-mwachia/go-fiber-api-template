package routes

import (
	"time"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/clinton-mwachia/go-fiber-api-template/controllers"
	"github.com/clinton-mwachia/go-fiber-api-template/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetUpRouter(app *fiber.App) {
	app.Use(logger.New())

	api := app.Group("/api")

	api.Post("/login", controllers.Login)

	//api.Use(middlewares.AuthRequired())

	// get collections
	todoCollection := config.GetCollection("todos")

	// users routes
	api.Post("/user/register", controllers.Register)
	api.Get("/users", controllers.GetAllUsers)
	api.Get("/user/:id", controllers.GetUserByID)
	api.Get("/users/paginated", controllers.GetPaginatedUsers)
	api.Put("/user/:id", controllers.UpdateUser)
	api.Delete("/user/:id", controllers.DeleteUser)
	api.Put("/change-password/:id", controllers.ChangePassword)
	api.Put("/reset-password/:id", controllers.ResetPassword)

	// todos routes
	api.Post("/todo", controllers.CreateTodo)
	api.Get("/todos", controllers.GetTodos)
	api.Delete("/todo/:id", middlewares.EnsureTodoOwner(todoCollection), controllers.DeleteTodo)
	api.Put("/todo/:id", controllers.UpdateTodo)
	api.Get("/todo/:id", controllers.GetTodoByID)
	api.Get("/todos/:userId/count", controllers.CountTodosByUserID)
	// the api only make 3 requests per minute
	api.Get("/todos/count", limiter.New(limiter.Config{
		Max:        3,               // max requests
		Expiration: 1 * time.Minute, // per minute
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try after 1 minute",
			})
		},
	}), controllers.CountTodos)
	api.Get("/todos/:userId", controllers.GetTodosByUserID)
}
