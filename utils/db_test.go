package utils

import (
	"os"
	"testing"
)

func TestConnectDB(t *testing.T) {
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	ConnectDB()
	if DB == nil {
		t.Error("DB should not be nil after ConnectDB()")
	}
}

func TestGetCollection(t *testing.T) {
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	ConnectDB()
	coll := GetCollection("users")
	if coll == nil {
		t.Error("Collection should not be nil")
	}
}
