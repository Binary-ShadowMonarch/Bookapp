// internal/auth/jwt.go

package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("replace-me-with-env-var!") // load from ENV in prod

// Define TTLs for each token type
const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour // 7 days
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// CreateToken now accepts a TTL to generate either an access or refresh token.
func CreateToken(email string, ttl time.Duration) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(jwtSecret)
}

// ParseToken remains the same, it just verifies any token.
func ParseToken(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return nil, err
	}
	return tok.Claims.(*Claims), nil
}
