package utils

import (
	"testing"
)

func TestGenerateAndParseJWT(t *testing.T) {
	token, err := GenerateJWT("12345", "user")
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("failed to parse JWT: %v", err)
	}
	if claims.UserID != "12345" || claims.Role != "user" {
		t.Errorf("unexpected claims: %+v", claims)
	}
}
