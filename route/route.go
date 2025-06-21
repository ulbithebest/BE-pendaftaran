package route

import (
	"github.com/gofiber/fiber/v2"
	"ulbithebest/BE-pendaftaran/controller"
	"ulbithebest/BE-pendaftaran/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth
	api.Post("/register", controller.RegisterUser)
	api.Post("/login", controller.LoginUser)

	// User routes (JWT required)
	user := api.Group("", middleware.JWTAuthMiddleware())
	user.Get("/me", controller.GetMe)
	user.Post("/registration", controller.SubmitRegistration)
	user.Get("/registration", controller.GetMyRegistration)
	user.Post("/upload/cv", controller.UploadCV)

	// Admin routes (JWT + admin only)
	admin := api.Group("/admin", middleware.JWTAuthMiddleware(), middleware.AdminOnlyMiddleware())
	admin.Get("/registrations", controller.ListRegistrations)
	admin.Get("/registration/:id", controller.GetRegistrationDetail)
	admin.Post("/registration/:id/status", controller.UpdateRegistrationStatus)
}
