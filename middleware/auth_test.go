package middleware

import (
	"net/http/httptest"
	"testing"
	"ulbithebest/BE-pendaftaran/utils"

	"github.com/gofiber/fiber/v2"
)

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	app := fiber.New()
	app.Use(JWTAuthMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	token, _ := utils.GenerateJWT("userid123", "user")
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(JWTAuthMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestJWTAuthMiddleware_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(JWTAuthMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAdminOnlyMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return c.Next()
	})
	app.Use(AdminOnlyMiddleware())
	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.SendString("admin ok")
	})
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAdminOnlyMiddleware_Forbidden(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "user")
		return c.Next()
	})
	app.Use(AdminOnlyMiddleware())
	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.SendString("admin ok")
	})
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Errorf("expected 403, got %d", resp.StatusCode)
	}
}
