package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Load JWT secret from env variables
var secret = os.Getenv("JWT_SECRET")

// Convert secret to byte slice for signing
var jwtSecret = []byte(secret)


func GenerateJWT(userID int64, role string) (string, error) {
	
	// Define token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"role"   : role,
		"exp"    : time.Now().Add(15 * time.Minute).Unix(), // Token expires in 15 minutes
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtSecret)
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


/* func GenerateSecureToken() string {
	b := make([]byte, 32) // Gen a 32-byte random token
	rand.Read(b)
	return base64.RawStdEncoding.EncodeToString(b)
}
*/