package utils

import (
	"testing"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
)

func TestSuccessResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/success", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.Map{"foo": "bar"})
	})
	req := httptest.NewRequest("GET", "/success", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/error", func(c *fiber.Ctx) error {
		return ErrorResponse(c, 400, "fail")
	})
	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
