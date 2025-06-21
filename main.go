package main

import (
	"github.com/gofiber/fiber/v2"
	"ulbithebest/BE-pendaftaran/config"
	"ulbithebest/BE-pendaftaran/route"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	route.SetupRoutes(app)
	app.Listen(":8080")
}
