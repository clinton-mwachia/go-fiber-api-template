package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/clinton-mwachia/go-fiber-api-template/controllers"
	"github.com/clinton-mwachia/go-fiber-api-template/routes"
	"github.com/clinton-mwachia/go-fiber-api-template/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// cors config for customization
	app.Use(cors.New(cors.Config{
		// user "*" in AllowOrigins to allow all origins, methods etc but it is prohibited
		// because it can expose your application to security risks.
		AllowOrigins: "https://example.com, https://example.com",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// ensure uploads folder is created
	utils.EnsureUploadsFolder()

	// load env
	config.Load()
	// connect DB
	config.ConnectDB()

	// initialize user collection
	controllers.InitUserCollection()
	controllers.InitTodoCollection()

	// setup routes (controllers contain logic)
	routes.SetUpRouter(app)

	// Handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// graceful shutdown
	go func() {
		if err := app.Listen(":" + config.Cfg.Port); err != nil {
			log.Printf("listen error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("🛑 Shutting down server...")

	// Disconnect DB gracefully
	config.DisconnectDB()
}
