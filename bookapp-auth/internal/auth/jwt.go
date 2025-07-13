package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("replace-me-with-env-var!") // load from ENV in prod
const jwtTTL = 24 * time.Hour

// Claims holds whatever you want inside the token.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// CreateToken signs a JWT for the given email.
func CreateToken(email string) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(jwtSecret)
}

// ParseToken verifies and returns claims, or an error.
func ParseToken(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return nil, err
	}
	return tok.Claims.(*Claims), nil
}
