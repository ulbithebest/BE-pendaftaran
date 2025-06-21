package controllers

import (
	"context"
	"time"
	"strings"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"ulbithebest/BE-pendaftaran/models"
	"ulbithebest/BE-pendaftaran/utils"
)

// RegisterUser handles user registration
func RegisterUser(c *fiber.Ctx) error {
	type reqBody struct {
		Name     string `json:"name"`
		NIM      string `json:"nim"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Name == "" || body.NIM == "" || body.Email == "" || body.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "All fields are required")
	}
	// Check if email or NIM already exists
	userCol := utils.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, _ := userCol.CountDocuments(ctx, bson.M{"$or": []bson.M{{"email": strings.ToLower(body.Email)}, {"nim": body.NIM}}})
	if count > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email or NIM already registered")
	}
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}
	user := models.User{
		ID:       primitive.NewObjectID(),
		Name:     body.Name,
		NIM:      body.NIM,
		Email:    strings.ToLower(body.Email),
		Password: string(hash),
		Role:     "user",
	}
	_, err = userCol.InsertOne(ctx, user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user")
	}
	return utils.SuccessResponse(c, fiber.Map{"message": "Registration successful"})
}

// LoginUser handles user login and returns JWT
func LoginUser(c *fiber.Ctx) error {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	userCol := utils.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCol.FindOne(ctx, bson.M{"email": strings.ToLower(body.Email)}).Decode(&user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}
	return utils.SuccessResponse(c, fiber.Map{"token": token, "role": user.Role})
}

// GetMe returns current user info
func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userCol := utils.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(userID)
	err := userCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	user.Password = ""
	return utils.SuccessResponse(c, user)
}
