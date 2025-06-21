package controllers

import (
	"testing"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
	"strings"
)

func TestRegisterUser(t *testing.T) {
	app := fiber.New()
	app.Post("/register", RegisterUser)
	body := `{"name": "Test", "nim": "1234567890", "email": "test@example.com", "password": "password"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		t.Errorf("expected status 200 or 400, got %d", resp.StatusCode)
	}
}

func TestLoginUser(t *testing.T) {
	app := fiber.New()
	app.Post("/login", LoginUser)
	body := `{"email": "test@example.com", "password": "password"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		t.Errorf("expected status 200 or 401, got %d", resp.StatusCode)
	}
}
