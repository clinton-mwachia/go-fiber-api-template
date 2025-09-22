package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/clinton-mwachia/go-fiber-api-template/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var todoCollection *mongo.Collection

// Init sets up the collections after DB connection
func InitTodoCollection() {
	todoCollection = config.GetCollection("todos")
}

func CreateTodo(c *fiber.Ctx) error {
	title := c.FormValue("title")
	userID := c.FormValue("userId")

	var user models.User

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// confirm user exists
	err = userCollection.FindOne(context.Background(), bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		fmt.Println(err.Error())
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	if title == "" || userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title & UserID is required"})
	}

	todo := models.Todo{
		ID:        primitive.NewObjectID(),
		UserID:    uid,
		Title:     title,
		Completed: false,
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		// Save file
		filename := fmt.Sprintf("uploads/%s", strings.ToLower(file.Filename))
		if err := c.SaveFile(file, filename); err != nil {
			fmt.Println(err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save image"})
		}
		todo.Image = filename
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create todo"})
	}

	res := models.Todo{
		ID:        todo.ID,
		UserID:    todo.UserID,
		Title:     todo.Title,
		Completed: todo.Completed,
		Image:     todo.Image,
	}

	return c.Status(201).JSON(res)
}
