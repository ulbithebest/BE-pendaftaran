package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua variabel konfigurasi aplikasi
type Config struct {
	MongoURI            string
	DatabaseName        string
	PasetoSecretKey     string
	ServerPort          string
	CloudinaryCloudName string
	CloudinaryApiKey    string
	CloudinaryApiSecret string
}

var appConfig *Config

// getEnvWithDefault mengambil environment variable dengan fallback ke default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig memuat konfigurasi dari file .env dan environment variables
func LoadConfig() {
	// Coba load .env file (untuk development lokal)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	// Load basic config yang diperlukan untuk koneksi database
	appConfig = &Config{
		MongoURI:     getEnvWithDefault("MONGO_URI", ""),
		DatabaseName: getEnvWithDefault("MONGO_DATABASE", "himatif"),
		ServerPort:   getEnvWithDefault("SERVER_PORT", ":8080"),
	}

	// Validasi konfigurasi penting untuk koneksi database
	if appConfig.MongoURI == "" {
		log.Fatal("MONGO_URI is required")
	}

	log.Println("✅ Basic configuration loaded successfully")
}

// LoadDatabaseCredentials memuat credentials dari database setelah koneksi terbentuk
func LoadDatabaseCredentials(credentials map[string]string) {
	if appConfig == nil {
		log.Fatal("Basic config must be loaded first")
	}

	// Load credentials dari database dengan fallback ke environment variables
	appConfig.PasetoSecretKey = getCredentialWithFallback(credentials, "PASETO_SECRET_KEY", "")
	appConfig.CloudinaryCloudName = getCredentialWithFallback(credentials, "CLOUDINARY_CLOUD_NAME", "")
	appConfig.CloudinaryApiKey = getCredentialWithFallback(credentials, "CLOUDINARY_API_KEY", "")
	appConfig.CloudinaryApiSecret = getCredentialWithFallback(credentials, "CLOUDINARY_API_SECRET", "")

	// Validasi credentials yang wajib ada
	if appConfig.PasetoSecretKey == "" {
		log.Fatal("PASETO_SECRET_KEY is required (not found in database or environment)")
	}

	log.Printf("✅ Database credentials loaded successfully")
	log.Printf("   - PASETO_SECRET_KEY: %s", maskCredential(appConfig.PasetoSecretKey))
	log.Printf("   - CLOUDINARY_CLOUD_NAME: %s", maskCredential(appConfig.CloudinaryCloudName))
	log.Printf("   - CLOUDINARY_API_KEY: %s", maskCredential(appConfig.CloudinaryApiKey))
}

// getCredentialWithFallback mengambil credential dari database, fallback ke env variable
func getCredentialWithFallback(credentials map[string]string, key, defaultValue string) string {
	// Coba ambil dari database terlebih dahulu
	if value, exists := credentials[key]; exists && value != "" {
		return value
	}

	// Fallback ke environment variable
	if envValue := os.Getenv(key); envValue != "" {
		log.Printf("Warning: Using environment variable for %s (not found in database)", key)
		return envValue
	}

	return defaultValue
}

// maskCredential untuk menyembunyikan sebagian credential dalam log
func maskCredential(credential string) string {
	if len(credential) <= 4 {
		return "****"
	}
	return credential[:4] + "****"
}

// GetConfig mengembalikan instance konfigurasi yang sudah dimuat
func GetConfig() *Config {
	if appConfig == nil {
		LoadConfig()
	}
	return appConfig
}
