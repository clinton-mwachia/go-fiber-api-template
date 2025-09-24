package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/clinton-mwachia/go-fiber-api-template/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// EnsureTodoOwner ensures that only the user who created the todo can update/delete it
func EnsureTodoOwner(todoCollection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get todo ID from URL
		todoIDParam := c.Params("id")
		todoID, err := primitive.ObjectIDFromHex(todoIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
		}

		// Get user ID from context (set by AuthRequired middleware)
		userIDStr := c.Locals("user_id").(string)
		userID, _ := primitive.ObjectIDFromHex(userIDStr)

		// Find the todo
		var todo models.Todo
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = todoCollection.FindOne(ctx, primitive.M{"_id": todoID}).Decode(&todo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Check ownership
		if todo.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not allowed to modify this todo"})
		}

		return c.Next()
	}
}
