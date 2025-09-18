package controllers

import (
	"context"
	"log"
	"time"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/clinton-mwachia/go-fiber-api-template/models"
	"github.com/clinton-mwachia/go-fiber-api-template/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	body.ID = primitive.NewObjectID()
	// set ID manually
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := userCollection.InsertOne(ctx, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

// get all users
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	cursor, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &users); err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse users"})
	}

	return c.JSON(users)
}
