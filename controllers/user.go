package controllers

import (
	"context"
	"time"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/clinton-mwachia/go-fiber-api-template/models"
	"github.com/clinton-mwachia/go-fiber-api-template/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// define user collection
var userCollection *mongo.Collection

// Init sets up the collections after DB connection
func InitUserCollection() {
	userCollection = config.GetCollection("users")
}

// register a new user
func Register(c *fiber.Ctx) error {
	var body models.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Hash password
	hashed, _ := utils.HashPassword(body.Password)
	body.Password = hashed
	if body.Role == "" {
		body.Role = "user"
	}

	// set ID manually
	// body.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := userCollection.InsertOne(ctx, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}
