package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"ulbithebest/BE-pendaftaran/utils"
)

// JWTAuthMiddleware checks JWT token validity
func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid token")
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseJWT(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

// AdminOnlyMiddleware restricts access to admin role
func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin access required")
		}
		return c.Next()
	}
}
