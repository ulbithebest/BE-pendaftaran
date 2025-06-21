package routes

import (
	"github.com/gofiber/fiber/v2"
	"ulbithebest/BE-pendaftaran/controllers"
	"ulbithebest/BE-pendaftaran/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth
	api.Post("/register", controllers.RegisterUser)
	api.Post("/login", controllers.LoginUser)

	// User routes (JWT required)
	user := api.Group("", middleware.JWTAuthMiddleware())
	user.Get("/me", controllers.GetMe)
	user.Post("/registration", controllers.SubmitRegistration)
	user.Get("/registration", controllers.GetMyRegistration)
	user.Post("/upload/cv", controllers.UploadCV)

	// Admin routes (JWT + admin only)
	admin := api.Group("/admin", middleware.JWTAuthMiddleware(), middleware.AdminOnlyMiddleware())
	admin.Get("/registrations", controllers.ListRegistrations)
	admin.Get("/registration/:id", controllers.GetRegistrationDetail)
	admin.Post("/registration/:id/status", controllers.UpdateRegistrationStatus)
}
