package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectDB membuat koneksi ke MongoDB Atlas
func ConnectDB(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	MongoClient = client
	log.Println("Successfully connected to MongoDB!")
}

// GetConfigCredentials mengambil semua credentials dari collection configurasi di database himatif
func GetConfigCredentials() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Access database 'himatif' and collection 'configurasi'
	collection := MongoClient.Database("himatif").Collection("configurasi")

	var config model.ConfigCredential
	err := collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no configuration document found in himatif.configurasi")
		}
		return nil, fmt.Errorf("failed to fetch config credentials: %v", err)
	}

	credentials := make(map[string]string)

	// Map struct fields to credential keys
	if config.CloudinaryAPIKey != "" {
		credentials["CLOUDINARY_API_KEY"] = config.CloudinaryAPIKey
	}
	if config.CloudinaryAPISecret != "" {
		credentials["CLOUDINARY_API_SECRET"] = config.CloudinaryAPISecret
	}
	if config.CloudinaryCloudName != "" {
		credentials["CLOUDINARY_CLOUD_NAME"] = config.CloudinaryCloudName
	}
	if config.PasetoSecretKey != "" {
		credentials["PASETO_SECRET_KEY"] = config.PasetoSecretKey
	}
	if config.ServerPort != "" {
		credentials["SERVER_PORT"] = config.ServerPort
	}

	log.Printf("✅ Loaded %d credentials from database", len(credentials))
	return credentials, nil
}

// GetConfigCredential mengambil satu credential berdasarkan key
func GetConfigCredential(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("himatif").Collection("configurasi")

	var config model.ConfigCredential
	err := collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("configuration document not found in database")
		}
		return "", fmt.Errorf("failed to fetch credential '%s': %v", key, err)
	}

	// Return the specific field based on key
	switch key {
	case "CLOUDINARY_API_KEY":
		return config.CloudinaryAPIKey, nil
	case "CLOUDINARY_API_SECRET":
		return config.CloudinaryAPISecret, nil
	case "CLOUDINARY_CLOUD_NAME":
		return config.CloudinaryCloudName, nil
	case "PASETO_SECRET_KEY":
		return config.PasetoSecretKey, nil
	case "SERVER_PORT":
		return config.ServerPort, nil
	default:
		return "", fmt.Errorf("credential '%s' not found in configuration", key)
	}
}
