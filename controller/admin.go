package controller

import (
	"context"
	"time"
	"ulbithebest/BE-pendaftaran/config"
	"ulbithebest/BE-pendaftaran/helper"
	"ulbithebest/BE-pendaftaran/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListRegistrations(c *fiber.Ctx) error {
	regCol := config.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := regCol.Find(ctx, bson.M{})
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch registrations")
	}
	var regs []models.Registration
	if err := cursor.All(ctx, &regs); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to decode registrations")
	}
	return helper.SuccessResponse(c, regs)
}

func GetRegistrationDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid registration ID")
	}
	regCol := config.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var reg models.Registration
	err = regCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&reg)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Registration not found")
	}
	return helper.SuccessResponse(c, reg)
}

func UpdateRegistrationStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid registration ID")
	}
	type reqBody struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Status != "lulus" && body.Status != "tidak_lulus" && body.Status != "menunggu" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status value")
	}
	regCol := config.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{"$set": bson.M{"status": body.Status, "note": body.Note, "updated_at": time.Now()}}
	res := regCol.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update)
	if res.Err() != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Registration not found or failed to update")
	}
	return helper.SuccessResponse(c, fiber.Map{"message": "Status updated"})
}
