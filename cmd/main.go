package main

import (
	"github.com/gofiber/fiber/v2"
	"ulbithebest/BE-pendaftaran/routes"
	"ulbithebest/BE-pendaftaran/utils"
)

func main() {
	// Connect to MongoDB
	utils.ConnectDB()

	app := fiber.New()

	// Register routes
	routes.SetupRoutes(app)

	// Start server
	app.Listen(":8080")
}
