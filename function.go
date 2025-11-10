package pendaftaran

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

    "github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Global router instance
var (
    router *chi.Mux
    once   sync.Once
)

// allowlist for CORS preflight handler (manual path)
var allowedOrigins = map[string]struct{}{
    "https://ulbithebest.github.io": {},
    "http://localhost:5500":         {},
    "http://127.0.0.1:5500":         {},
    "http://127.0.0.1:5501":         {},
    "http://localhost:5501":         {},
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
    origin := r.Header.Get("Origin")
    if _, ok := allowedOrigins[origin]; !ok {
        return
    }
    // Mirror the request origin (required when credentials are used)
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Vary", "Origin")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Requested-With")
    w.Header().Set("Access-Control-Max-Age", "300")
}

// initializeApp initializes the application router and database connection
func initializeApp() {
	log.Println("🚀 Initializing Cloud Function...")

	// 1. Load configuration (MONGO_URI, MONGO_DATABASE, etc)
	cfg := config.GetConfig()
	log.Printf("✅ Config loaded - DB: %s, Port: %s", cfg.DatabaseName, cfg.ServerPort)

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
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)

	// 6. Setup CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{
			"https://ulbithebest.github.io", // GitHub Pages frontend
			"http://localhost:5500",         // Local dev
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

	// 7. Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"pendaftaran-api"}`))
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

		// File server (protected)
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

	router = r
	log.Println("✅ Router initialized successfully")
}

// URL handles all HTTP requests - entry point for Cloud Function
func URL(w http.ResponseWriter, r *http.Request) {
    // Initialize only once
    once.Do(initializeApp)

    // Handle preflight (OPTIONS) manually to ensure Cloud Functions compatibility
    if r.Method == http.MethodOptions {
        setCORSHeaders(w, r)
        w.WriteHeader(http.StatusNoContent)
        return
    }

	// Ensure router exists
	if router == nil {
		log.Println("❌ Router not initialized")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Forward request to Chi router
	router.ServeHTTP(w, r)
}

func init() {
	functions.HTTP("Pendaftaran", URL)
}
