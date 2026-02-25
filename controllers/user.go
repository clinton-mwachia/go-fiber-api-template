package controllers

import (
	"context"
	"os"
	"time"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/clinton-mwachia/go-fiber-api-template/models"
	"github.com/clinton-mwachia/go-fiber-api-template/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// define user collection
var userCollection *mongo.Collection

// Login request body
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login response with expiry
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // UNIX timestamp
}

// Init sets up the collections after DB connection
func InitUserCollection() {
	userCollection = config.GetCollection("users")
}

// register a new user
func Register(c *fiber.Ctx) error {
	var body models.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

// get all users
func GetAllUsers(c *fiber.Ctx) error {
	cursor, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users: " + err.Error()})
	}
	defer cursor.Close(context.Background())

	users := []models.User{}
	if err := cursor.All(context.Background(), &users); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse users"})
	}

	return c.JSON(users)
}

// get all users with pagination
func GetPaginatedUsers(c *fiber.Ctx) error {
	var users []models.User

	// pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}) // newest first

	cursor, err := userCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users: " + err.Error()})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse users: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"page":  page,
		"limit": limit,
		"data":  users,
	})
}

// get user by id
func GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	var body struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Role     *string `json:"role"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user: " + err.Error()})
	}
	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Fetch updated user
	var updatedUser models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch updated user: " + err.Error()})
	}

	return c.JSON(updatedUser)
}

// delete user by id
func DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user: " + err.Error()})
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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password: " + err.Error()})
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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}

	userId := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reset password: " + err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

// Login user
func Login(c *fiber.Ctx) error {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}

	// Find user by email
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return c.Status(404).JSON(fiber.Map{"error": "User not found: " + err.Error()})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Something went wrong: " + err.Error()})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password: " + err.Error()})
	}

	// Expiry time
	expirationTime := time.Now().Add(time.Hour * 72).Unix() // 72 hours from now
	// Create JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     expirationTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(LoginResponse{
		Token:     signedToken,
		ExpiresAt: expirationTime,
	})
}
