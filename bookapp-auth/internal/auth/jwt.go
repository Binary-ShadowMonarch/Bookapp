// internal/auth/jwt.go

package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Define TTLs for each token type
const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour // 7 days
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func getJwtSecret() []byte {
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		log.Fatal("AUTH_JWT_SECRET is not set")
	}
	return []byte(secret)
}

// CreateToken now accepts a TTL to generate either an access or refresh token.
func CreateToken(email string, ttl time.Duration) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // <-- unique per token
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ttl)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(getJwtSecret())
}

// ParseToken remains the same, it just verifies any token.
func ParseToken(tokenStr string) (*Claims, error) {

	log.Printf("env : %s", getJwtSecret())
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if err != nil || !tok.Valid {
		return nil, err
	}
	return tok.Claims.(*Claims), nil
}
