// Package main defines the Cloud Function entry point
// This file is designed for Google Cloud Functions Gen 2 with Functions Framework
package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/handler"
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Global router instance
var (
	router *chi.Mux
	once   sync.Once
)

// initializeApp initializes the application router and database connection
func initializeApp() {
	log.Println("🚀 Initializing Cloud Function...")

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

	// 6. Setup global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)

	// 7. Setup CORS - configured for production
	corsOptions := cors.Options{
		AllowedOrigins: []string{
			"https://svalvva.github.io", // GitHub Pages frontend
			"http://localhost:5500",     // Local development
			"http://127.0.0.1:5500",
			"http://127.0.0.1:5501",
			"http://localhost:5501",
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.Handler(corsOptions))

	// 6. Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"pendaftaran-api"}`))
	})

	// 7. Public routes (no authentication required)
	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	// 8. Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// User endpoints
		r.Get("/user/profile", handler.GetUserProfileHandler)
		r.Post("/user/registration", handler.SubmitRegistrationHandler)
		r.Get("/user/my-registration", handler.GetUserRegistrationHandler)
		r.Get("/info", handler.GetAllInfoHandler)

		// File server for uploads (protected)
		if _, err := os.Stat("./uploads"); err == nil {
			fileServer := http.FileServer(http.Dir("./uploads"))
			r.Handle("/uploads/*", http.StripPrefix("/api/uploads/", fileServer))
		}

		// Admin-only routes
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

	// Assign to global variable
	router = r
	log.Println("✅ Router initialized successfully")
}

// Pendaftaran is the main Cloud Function entry point
// This function name must match the --entry-point in gcloud deploy
func Pendaftaran(w http.ResponseWriter, r *http.Request) {
	// Initialize app only once using sync.Once
	once.Do(initializeApp)

	// Handle the request
	if router == nil {
		log.Println("❌ Router not initialized")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Serve the request
	router.ServeHTTP(w, r)
}

func init() {
	// Register the Cloud Function
	funcframework.RegisterHTTPFunction("/", Pendaftaran)
}

// main function for Functions Framework
func main() {
	// Use PORT environment variable, or default to 8080
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
