package controllers

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"
	"ulbithebest/BE-pendaftaran/models"
	"ulbithebest/BE-pendaftaran/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestListRegistrations(t *testing.T) {
	app := fiber.New()
	app.Get("/admin/registrations", ListRegistrations)

	// Insert dummy data
	regCol := utils.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	regCol.InsertOne(ctx, models.Registration{
		ID:         primitive.NewObjectID(),
		Division:   "TestDiv",
		Motivation: "TestMotivation",
		Status:     "menunggu",
		UpdatedAt:  time.Now(),
	})

	req := httptest.NewRequest("GET", "/admin/registrations", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetRegistrationDetail(t *testing.T) {
	app := fiber.New()
	app.Get("/admin/registration/:id", GetRegistrationDetail)

	regCol := utils.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reg := models.Registration{
		ID:         primitive.NewObjectID(),
		Division:   "TestDiv",
		Motivation: "TestMotivation",
		Status:     "menunggu",
		UpdatedAt:  time.Now(),
	}
	regCol.InsertOne(ctx, reg)

	req := httptest.NewRequest("GET", "/admin/registration/"+reg.ID.Hex(), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
