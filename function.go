package main // or another package name like 'pendaftaran'

import (
	"net/http"

	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/handler"
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Declare the router as a global variable so it's initialized only once.
var router *chi.Mux

func init() {
	// This init() function runs once when the Cloud Function instance starts.

	// 1. Muat Konfigurasi dari file .env
	cfg := config.GetConfig()

	// 2. Hubungkan ke Database MongoDB
	repository.ConnectDB(cfg)

	// 3. Inisialisasi Router menggunakan Chi
	r := chi.NewRouter()

	// 4. Setup Middleware Global
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Setup CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Use "*" for broader access or specify your frontend URLs
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 5. Definisikan Routes (Endpoint API)
	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/user/profile", handler.GetUserProfileHandler)
		r.Post("/user/registration", handler.SubmitRegistrationHandler)
		r.Get("/user/my-registration", handler.GetUserRegistrationHandler)
		r.Get("/info", handler.GetAllInfoHandler)
		
		fileServer := http.FileServer(http.Dir("./uploads"))
		r.Handle("/uploads/*", http.StripPrefix("/api/uploads/", fileServer))

		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminOnlyMiddleware)
			r.Get("/registrations-with-details", handler.GetAllRegistrationsDetailHandler)
			r.Patch("/registrations/{id}", handler.UpdateRegistrationDetailsHandler)
			r.Patch("/registrations/bulk-update", handler.BulkUpdateStatusHandler)
			r.Get("/users", handler.GetAllUsersHandler)
			r.Delete("/registrations/{id}", handler.DeleteRegistrationHandler)
			r.Post("/info", handler.CreateInfoHandler)
			r.Put("/info/{id}", handler.UpdateInfoHandler)
			r.Delete("/info/{id}", handler.DeleteInfoHandler)
		})
	})

	// Assign the configured router to the global variable
	router = r
}

// Pendaftaran is the exportable handler function for Google Cloud Functions.
// The name MUST be capitalized to be exported.
func Pendaftaran(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}