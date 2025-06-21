package config

import "os"

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "supersecretkey"
	}
	return secret
}
