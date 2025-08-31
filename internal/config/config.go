package config

import (
	"log"
	"os"
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

// LoadConfig memuat konfigurasi dari environment variables sistem
func LoadConfig() {
	log.Println("🔧 Loading configuration from system environment variables...")

	// Debug: print semua environment variables yang relevan
	mongoURI := getEnvWithDefault("MONGO_URI", "")
	mongoDatabase := getEnvWithDefault("MONGO_DATABASE", "pendaftaran_db")
	serverPort := getEnvWithDefault("SERVER_PORT", ":8080")

	log.Printf("🔍 Debug - MONGO_URI length: %d", len(mongoURI))
	log.Printf("🔍 Debug - MONGO_DATABASE: %s", mongoDatabase)
	log.Printf("🔍 Debug - SERVER_PORT: %s", serverPort)

	// Load basic config yang diperlukan untuk koneksi database
	appConfig = &Config{
		MongoURI:     mongoURI,
		DatabaseName: mongoDatabase,
		ServerPort:   serverPort,
	}

	// Validasi konfigurasi penting untuk koneksi database
	if appConfig.MongoURI == "" {
		log.Fatal("❌ MONGO_URI is required - please set environment variable MONGO_URI")
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
		log.Printf("✅ Using credential %s from database", key)
		return value
	}
	
	// Fallback ke environment variable (untuk backward compatibility)
	if envValue := os.Getenv(key); envValue != "" {
		log.Printf("⚠️ Warning: Using environment variable for %s (not found in database)", key)
		return envValue
	}
	
	log.Printf("❌ Credential %s not found in database or environment", key)
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