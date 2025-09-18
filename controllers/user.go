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
	"golang.org/x/crypto/bcrypt"
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

// get user by id
func GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	return c.JSON(user)
}

// update user by id
func UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var body struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Role     *string `json:"role"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	update := bson.M{}
	if body.Username != nil {
		update["username"] = body.Username
	}
	if body.Email != nil {
		update["email"] = *body.Email
	}
	if body.Role != nil {
		update["role"] = *body.Role
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No fields to update"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}
	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Fetch updated user
	var updatedUser models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch updated user"})
	}

	return c.JSON(updatedUser)
}

// delete user by id
func DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// change password
func ChangePassword(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.CurrentPassword == "" || body.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Both current and new password are required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch user
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	// Verify current password
	if !utils.CheckPassword(user.Password, body.CurrentPassword) {
		return c.Status(400).JSON(fiber.Map{"error": "Current password is incorrect"})
	}

	// Hash new password
	hashed, _ := utils.HashPassword(body.NewPassword)

	// Update in DB
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"password": hashed}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
	}

	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}

// reset user password
// ONLY ADMIN CAN DO THIS
func ResetPassword(c *fiber.Ctx) error {
	type ResetInput struct {
		NewPassword string `json:"newPassword"`
	}

	var input ResetInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userId := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Update the user’s password
	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}
	result, err := userCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}
