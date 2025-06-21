package controllers

import (
	"context"
	"fmt"
	// "mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ulbithebest/BE-pendaftaran/models"
	"ulbithebest/BE-pendaftaran/utils"
)

// SubmitRegistration handles user registration form submission
func SubmitRegistration(c *fiber.Ctx) error {
	type reqBody struct {
		Division   string `json:"division"`
		Motivation string `json:"motivation"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Division == "" || body.Motivation == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Division and motivation are required")
	}
	userID := c.Locals("user_id").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)
	regCol := utils.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Check if already registered
	count, _ := regCol.CountDocuments(ctx, bson.M{"user_id": objID})
	if count > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "You have already registered")
	}
	reg := models.Registration{
		ID:         primitive.NewObjectID(),
		UserID:     objID,
		Division:   body.Division,
		Motivation: body.Motivation,
		CVPath:     "",
		Status:     "menunggu",
		Note:       "",
		UpdatedAt:  time.Now(),
	}
	_, err := regCol.InsertOne(ctx, reg)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to submit registration")
	}
	return utils.SuccessResponse(c, fiber.Map{"message": "Registration submitted"})
}

// GetMyRegistration returns user's registration data
func GetMyRegistration(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)
	regCol := utils.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var reg models.Registration
	err := regCol.FindOne(ctx, bson.M{"user_id": objID}).Decode(&reg)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No registration found")
	}
	return utils.SuccessResponse(c, reg)
}

// UploadCV handles CV file upload
func UploadCV(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)
	file, err := c.FormFile("cv")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "CV file is required")
	}
	if !strings.HasSuffix(file.Filename, ".pdf") {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "CV must be a PDF file")
	}
	// Save file
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
	filename := fmt.Sprintf("cv_%s_%d.pdf", userID, time.Now().Unix())
	path := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, path); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save CV file")
	}
	// Update registration
	regCol := utils.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := regCol.FindOneAndUpdate(ctx, bson.M{"user_id": objID}, bson.M{"$set": bson.M{"cv_path": path, "updated_at": time.Now()}})
	if res.Err() != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Registration not found or failed to update CV")
	}
	return utils.SuccessResponse(c, fiber.Map{"cv_path": path})
}
