package unit

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmashedxtomatox/auth-service/internal/utils"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	// Arrange
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		t.Fatal("JWT_SECRET environment variable is not set")
	}

	userID := 123

	// Act
	token, err := utils.GenerateJWT(userID, jwtSecret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims := &utils.Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	// Assert
	if err != nil || !parsed.Valid {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Errorf("expected token expiry in the future, got %v", claims.ExpiresAt)
	}
}
