package controllers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"ulbithebest/BE-pendaftaran/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func TestSubmitRegistration(t *testing.T) {
	app := fiber.New()
	app.Post("/registration", SubmitRegistration)
	regCol := utils.GetCollection("registrations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	regCol.DeleteMany(ctx, bson.M{"division": "TestDiv"})

	body := `{"division": "TestDiv", "motivation": "TestMotivation"}`
	req := httptest.NewRequest("POST", "/registration", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// Simulasikan JWT auth dengan context (bisa diimprove dengan middleware mock)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		t.Errorf("expected status 200 or 400, got %d", resp.StatusCode)
	}
}
