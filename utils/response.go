package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Success response
func SuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// Error response
func ErrorResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
