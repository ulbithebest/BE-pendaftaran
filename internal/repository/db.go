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
		log.Printf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return
	}

	MongoClient = client
	log.Println("Successfully connected to MongoDB!")
}

// GetConfigCredentials mengambil semua credentials dari collection configurasi.
func GetConfigCredentials() (map[string]string, error) {
	if MongoClient == nil {
		return nil, fmt.Errorf("MongoDB client is not connected")
	}

	dbName := config.GetConfig().DatabaseName
	if dbName == "" {
		dbName = "himatif"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MongoClient.Database(dbName).Collection("configurasi")

	var credentialDoc model.ConfigCredential
	err := collection.FindOne(ctx, bson.M{}).Decode(&credentialDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no configuration document found in %s.configurasi", dbName)
		}
		return nil, fmt.Errorf("failed to fetch config credentials: %v", err)
	}

	credentials := make(map[string]string)

	// Map struct fields to credential keys
	if credentialDoc.CloudinaryAPIKey != "" {
		credentials["CLOUDINARY_API_KEY"] = credentialDoc.CloudinaryAPIKey
	}
	if credentialDoc.CloudinaryAPISecret != "" {
		credentials["CLOUDINARY_API_SECRET"] = credentialDoc.CloudinaryAPISecret
	}
	if credentialDoc.CloudinaryCloudName != "" {
		credentials["CLOUDINARY_CLOUD_NAME"] = credentialDoc.CloudinaryCloudName
	}
	if credentialDoc.PasetoSecretKey != "" {
		credentials["PASETO_SECRET_KEY"] = credentialDoc.PasetoSecretKey
	}
	if credentialDoc.ServerPort != "" {
		credentials["SERVER_PORT"] = credentialDoc.ServerPort
	}

	log.Printf("✅ Loaded %d credentials from database", len(credentials))
	return credentials, nil
}

// GetConfigCredential mengambil satu credential berdasarkan key
func GetConfigCredential(key string) (string, error) {
	if MongoClient == nil {
		return "", fmt.Errorf("MongoDB client is not connected")
	}

	dbName := config.GetConfig().DatabaseName
	if dbName == "" {
		dbName = "himatif"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database(dbName).Collection("configurasi")

	var credentialDoc model.ConfigCredential
	err := collection.FindOne(ctx, bson.M{}).Decode(&credentialDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("configuration document not found in database")
		}
		return "", fmt.Errorf("failed to fetch credential '%s': %v", key, err)
	}

	// Return the specific field based on key
	switch key {
	case "CLOUDINARY_API_KEY":
		return credentialDoc.CloudinaryAPIKey, nil
	case "CLOUDINARY_API_SECRET":
		return credentialDoc.CloudinaryAPISecret, nil
	case "CLOUDINARY_CLOUD_NAME":
		return credentialDoc.CloudinaryCloudName, nil
	case "PASETO_SECRET_KEY":
		return credentialDoc.PasetoSecretKey, nil
	case "SERVER_PORT":
		return credentialDoc.ServerPort, nil
	default:
		return "", fmt.Errorf("credential '%s' not found in configuration", key)
	}
}
