// internal/auth/jwt.go

package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// these are how long my tokens last
// access tokens are short-lived for security, refresh tokens last longer
const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour // 7 days
)

// Claims contains the data I store in my JWT tokens
// right now I just store the email, but I might add more later
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// getJwtSecret gets the secret key from environment variables
// this secret is used to sign and verify JWT tokens
func getJwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("AUTH_JWT_SECRET is not set")
	}
	return []byte(secret)
}

// CreateToken creates a new JWT token for a user
// I use this to create both access tokens and refresh tokens
// the ttl parameter determines how long the token lasts
func CreateToken(email string, ttl time.Duration) (string, error) {
	log.Printf("DEBUG: Creating JWT token for %s with TTL %v", email, ttl)
	
	// create the claims (data) that will be stored in the token
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // unique ID for each token
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ttl)),
		},
	}
	
	// create and sign the token with my secret key
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tok.SignedString(getJwtSecret())
	
	if err != nil {
		log.Printf("DEBUG: Failed to create JWT token for %s: %v", email, err)
	} else {
		log.Printf("DEBUG: Successfully created JWT token for %s", email)
	}
	
	return tokenString, err
}

// ParseToken verifies and parses a JWT token
// I use this to check if a token is valid and get the user's email from it
func ParseToken(tokenStr string) (*Claims, error) {
	log.Printf("DEBUG: Parsing JWT token")
	
	// parse the token and verify it's signed with my secret
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	
	if err != nil || !tok.Valid {
		log.Printf("DEBUG: JWT token parsing failed: %v", err)
		return nil, err
	}
	
	// extract the claims (user data) from the token
	claims := tok.Claims.(*Claims)
	log.Printf("DEBUG: Successfully parsed JWT token for %s", claims.Email)
	return claims, nil
}
