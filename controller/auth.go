package controller

import (
	"context"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ulbithebest/BE-pendaftaran/models"
	"ulbithebest/BE-pendaftaran/helper"
	"ulbithebest/BE-pendaftaran/config"
)

func RegisterUser(c *fiber.Ctx) error {
	type reqBody struct {
		Name     string `json:"name"`
		NIM      string `json:"nim"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Name == "" || body.NIM == "" || body.Email == "" || body.Password == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "All fields are required")
	}
	userCol := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, _ := userCol.CountDocuments(ctx, bson.M{"$or": []bson.M{{"email": strings.ToLower(body.Email)}, {"nim": body.NIM}}})
	if count > 0 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Email or NIM already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
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
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user")
	}
	return helper.SuccessResponse(c, fiber.Map{"message": "Registration successful"})
}

func LoginUser(c *fiber.Ctx) error {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	userCol := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCol.FindOne(ctx, bson.M{"email": strings.ToLower(body.Email)}).Decode(&user)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}
	token, err := helper.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}
	return helper.SuccessResponse(c, fiber.Map{"token": token, "role": user.Role})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userCol := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(userID)
	err := userCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	user.Password = ""
	return helper.SuccessResponse(c, user)
}
