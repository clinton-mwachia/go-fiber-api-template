package routes

import (
	"github.com/clinton-mwachia/go-fiber-api-template/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetUpRouter(app *fiber.App) {
	app.Use(logger.New())

	api := app.Group("/api")

	// users routes
	api.Post("/user/register", controllers.Register)
}
