package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "Admin@123456"
)

func HashPassword(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

func GenerateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func ResolveAdminUsername() string {
	value := os.Getenv("BILL_ADMIN_USERNAME")
	if value == "" {
		return defaultAdminUsername
	}

	return value
}

func ResolveAdminPasswordHash() string {
	value := os.Getenv("BILL_ADMIN_PASSWORD_SHA256")
	if value == "" {
		return HashPassword(defaultAdminPassword)
	}

	return value
}

func BuildAdminAccessToken(username string, passwordHash string) string {
	return HashPassword(username + ":" + passwordHash)
}

