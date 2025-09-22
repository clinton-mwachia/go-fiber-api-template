package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
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

// add a new todo
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
		filename := fmt.Sprintf("uploads/%s_%s", time.Now().Format("20060102150405"), strings.ToLower(file.Filename))
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

// get all todos
func GetTodos(c *fiber.Ctx) error {
	cursor, err := todoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch todos"})
	}
	defer cursor.Close(context.Background())

	todos := []models.Todo{}
	if err := cursor.All(context.Background(), &todos); err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse todos"})
	}

	return c.JSON(todos)
}

// delete todo by id
func DeleteTodo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	todoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Find the todo first to get image path
	var todo models.Todo
	err = todoCollection.FindOne(context.Background(), bson.M{"_id": todoID}).Decode(&todo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	// Delete the todo
	result, err := todoCollection.DeleteOne(context.Background(), bson.M{"_id": todoID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
	}
	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	// Delete image file if it exists
	if todo.Image != "" {
		if err := os.Remove(todo.Image); err != nil {
			c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("%s%s", "Failed to delete todo image: ", err.Error())})
		}
	}

	return c.JSON(fiber.Map{"message": "Todo deleted successfully"})
}

// update a todo
func UpdateTodo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	todoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Fetch current todo
	var todo models.Todo
	if err := todoCollection.FindOne(context.Background(), bson.M{"_id": todoID}).Decode(&todo); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	update := bson.M{}
	if body.Title != nil {
		update["title"] = *body.Title
	}
	if body.Completed != nil {
		update["completed"] = *body.Completed
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		// Delete old image if exists
		if todo.Image != "" {
			if err := os.Remove(todo.Image); err != nil {
				c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("%s%s", "Failed to delete old todo image:", err)})
			}
		}

		// Save new image
		filename := fmt.Sprintf("uploads/%s_%s", time.Now().Format("20060102150405"), strings.ToLower(file.Filename))
		if err := c.SaveFile(file, filename); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save new image"})
		}
		update["image"] = filename
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Nothing to update"})
	}

	// Update in MongoDB
	_, err = todoCollection.UpdateOne(context.Background(), bson.M{"_id": todoID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
	}

	// Return updated todo
	var updated models.Todo
	_ = todoCollection.FindOne(context.Background(), bson.M{"_id": todoID}).Decode(&updated)

	return c.JSON(updated)
}

// get todo by id
func GetTodoByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	todoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var todo models.Todo
	err = todoCollection.FindOne(context.Background(), bson.M{"_id": todoID}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch todo"})
	}

	return c.JSON(todo)
}
