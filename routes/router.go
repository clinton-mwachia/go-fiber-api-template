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
	api.Get("/users", controllers.GetAllUsers)
	api.Get("/user/:id", controllers.GetUserByID)
	api.Put("/user/:id", controllers.UpdateUser)
	api.Delete("/user/:id", controllers.DeleteUser)
	api.Put("/change-password/:id", controllers.ChangePassword)
	api.Put("/reset-password/:id", controllers.ResetPassword)
}
