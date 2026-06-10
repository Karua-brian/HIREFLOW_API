package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Load JWT secret from env variables
func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// GenerateJWT creates a JWT token with user ID and role as claims,
// signed with the secret key.
func GenerateJWT(userID uuid.UUID, role string) (string, error) {

	//	Define token claims
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Token expires in 15 minutes
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateRefreshToken() (string, error) {

	// Generate a random 32-byte token
	b := make([]byte, 32)

	// Fill the byte slice with random data
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a URL-safe base64 string
	return base64.URLEncoding.EncodeToString(b), nil
}
