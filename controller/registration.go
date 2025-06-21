package controller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ulbithebest/BE-pendaftaran/models"
	"ulbithebest/BE-pendaftaran/helper"
	"ulbithebest/BE-pendaftaran/config"
)

func SubmitRegistration(c *fiber.Ctx) error {
	type reqBody struct {
		Division   string `json:"division"`
		Motivation string `json:"motivation"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Division == "" || body.Motivation == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Division and motivation are required")
	}
	userID := c.Locals("user_id").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)
	regCol := config.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, _ := regCol.CountDocuments(ctx, bson.M{"user_id": objID})
	if count > 0 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "You have already registered")
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
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to submit registration")
	}
	return helper.SuccessResponse(c, fiber.Map{"message": "Registration submitted"})
}

func GetMyRegistration(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)
	regCol := config.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var reg models.Registration
	err := regCol.FindOne(ctx, bson.M{"user_id": objID}).Decode(&reg)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "No registration found")
	}
	return helper.SuccessResponse(c, reg)
}

func UploadCV(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	objID, _ := primitive.ObjectIDFromHex(userID)
	file, err := c.FormFile("cv")
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "CV file is required")
	}
	if !strings.HasSuffix(file.Filename, ".pdf") {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "CV must be a PDF file")
	}
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
	filename := fmt.Sprintf("cv_%s_%d.pdf", userID, time.Now().Unix())
	path := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, path); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save CV file")
	}
	regCol := config.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := regCol.FindOneAndUpdate(ctx, bson.M{"user_id": objID}, bson.M{"$set": bson.M{"cv_path": path, "updated_at": time.Now()}})
	if res.Err() != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Registration not found or failed to update CV")
	}
	return helper.SuccessResponse(c, fiber.Map{"cv_path": path})
}
