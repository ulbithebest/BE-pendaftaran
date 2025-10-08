package main

import (
	"log"
	"net/http"

	// "strings"

	// Pastikan path import ini sesuai dengan nama modul di go.mod Anda
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/handler"
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// 1. Load basic configuration (MONGO_URI, MONGO_DATABASE, SERVER_PORT)
	cfg := config.GetConfig()
	log.Printf("✅ Basic config loaded - DB: %s, Port: %s", cfg.DatabaseName, cfg.ServerPort)

	// 2. Connect to MongoDB
	repository.ConnectDB(cfg)
	log.Println("✅ Database connected")

	// 3. Load credentials from database himatif.configurasi
	credentials, err := repository.GetConfigCredentials()
	if err != nil {
		log.Printf("⚠️ Warning: Failed to load credentials from database: %v", err)
		log.Println("Will proceed with environment variables as fallback")
		credentials = make(map[string]string) // Empty map for fallback
	}

	// 4. Load database credentials into config
	config.LoadDatabaseCredentials(credentials)

	// 5. Initialize Chi router
	r := chi.NewRouter()

	// 6. Setup Middleware Global
	r.Use(chiMiddleware.Logger)    // Middleware untuk mencatat (log) setiap request yang masuk
	r.Use(chiMiddleware.Recoverer) // Middleware untuk menangani panic dan menjaga server tetap hidup

	// 7. Setup CORS (Cross-Origin Resource Sharing)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://ulbithebest.github.io", // GitHub Pages frontend
			"http://localhost:5500",         // Local development
			"http://127.0.0.1:5500",
			"http://127.0.0.1:5501",
			"http://localhost:5501",
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 8. Define Routes (API Endpoints)

	// Routes publik yang bisa diakses tanpa login/token
	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	// Group routes yang memerlukan otentikasi (wajib ada token Paseto)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Semua di dalam grup ini akan dilindungi oleh middleware otentikasi

		// --- Routes untuk user biasa ---
		r.Get("/user/profile", handler.GetUserProfileHandler)
		r.Post("/user/registration", handler.SubmitRegistrationHandler)
		r.Get("/user/my-registration", handler.GetUserRegistrationHandler)
		r.Get("/info", handler.GetAllInfoHandler)

		// --- TAMBAHAN: File Server untuk CV (dilindungi otentikasi) ---
		// Ini akan membuat file di folder /uploads bisa diakses via URL
		fileServer := http.FileServer(http.Dir("./uploads"))
		r.Handle("/uploads/*", http.StripPrefix("/api/uploads/", fileServer))

		// --- Routes khusus admin ---
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminOnlyMiddleware) // Perlindungan tambahan, hanya untuk admin

			// --- PERUBAHAN ENDPOINT ADMIN ---
			r.Get("/registrations-with-details", handler.GetAllRegistrationsDetailHandler)
			// r.Patch("/registrations/{id}/status", handler.UpdateRegistrationStatusHandler)
			r.Patch("/registrations/{id}", handler.UpdateRegistrationDetailsHandler)
			r.Patch("/registrations/bulk-update", handler.BulkUpdateStatusHandler)
			r.Get("/users", handler.GetAllUsersHandler)
			r.Delete("/registrations/{id}", handler.DeleteRegistrationHandler)
			r.Post("/info", handler.CreateInfoHandler)
			r.Put("/info/{id}", handler.UpdateInfoHandler)
			r.Delete("/info/{id}", handler.DeleteInfoHandler)
		})
	})

	// 9. Start HTTP Server
	log.Printf("✅ Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, r); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}