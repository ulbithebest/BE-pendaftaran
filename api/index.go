package handler

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/handler"
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	router *chi.Mux
	once   sync.Once
)

func initializeApp() {
	log.Println("🚀 Initializing Vercel Serverless Function...")

	// 1. Load configuration (MONGO_URI, MONGO_DATABASE, etc)
	cfg := config.GetConfig()
	log.Printf("✅ Config loaded - DB: %s", cfg.DatabaseName)

	// 2. Connect to MongoDB
	repository.ConnectDB(cfg)
	log.Println("✅ Database connected")

	// 3. Load credentials from database
	credentials, err := repository.GetConfigCredentials()
	if err != nil {
		log.Printf("⚠️ Failed to load credentials from DB: %v", err)
		log.Println("Fallback to environment variables")
		credentials = make(map[string]string)
	}
	config.LoadDatabaseCredentials(credentials)

	// 4. Setup Chi router
	r := chi.NewRouter()

	// 5. Global middlewares
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestLogMiddleware)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// 6. Setup CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{
			"https://ulbithebest.github.io", // GitHub Pages frontend
			"http://localhost:5500",         // Local dev
			"http://127.0.0.1:5500",
			"http://127.0.0.1:5501",
			"http://localhost:5501",
			"https://*.vercel.app",          // Vercel previews/deployments
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.Handler(corsOptions))

	// 7. Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"pendaftaran-api-vercel"}`))
	})

	// 8. Public routes
	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	// 9. Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// User endpoints
		r.Get("/user/profile", handler.GetUserProfileHandler)
		r.Post("/user/registration", handler.SubmitRegistrationHandler)
		r.Get("/user/my-registration", handler.GetUserRegistrationHandler)
		r.Get("/info", handler.GetAllInfoHandler)

		// File server (protected) - Note: Vercel serverless has ephemeral filesystem
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
			r.With(middleware.SuperAdminOnlyMiddleware).Patch("/users/{id}", handler.UpdateUserHandler)
			r.With(middleware.SuperAdminOnlyMiddleware).Patch("/users/{id}/password", handler.ResetUserPasswordHandler)
			r.Delete("/registrations/{id}", handler.DeleteRegistrationHandler)
			r.Post("/info", handler.CreateInfoHandler)
			r.Put("/info/{id}", handler.UpdateInfoHandler)
			r.Delete("/info/{id}", handler.DeleteInfoHandler)
			r.With(middleware.SuperAdminOnlyMiddleware).Get("/logs", handler.GetAppLogsHandler)
		})
	})

	router = r
	log.Println("✅ Router initialized successfully for Vercel")
}

// Handler is the entrypoint for Vercel Serverless Function
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initializeApp)

	if router == nil {
		log.Println("❌ Router not initialized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}

	router.ServeHTTP(w, r)
}
