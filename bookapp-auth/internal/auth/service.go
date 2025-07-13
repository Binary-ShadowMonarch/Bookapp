// internal/auth/service.go
package auth

import (
	"errors"
	"net/http"
	"strings"

	"bookapp/internal/models"
	// for CreateToken, ParseToken
)

// ErrInvalidCredentials is returned when login fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrUnauthorized is returned when a request lacks a valid JWT.
var ErrUnauthorized = errors.New("unauthorized")

// UserStore defines the methods we need for your store.
type UserStore interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	Update(*models.User) error
}

// Service wraps auth logic.
type Service struct {
	store UserStore
}

// NewService constructs the auth Service.
func NewService(us UserStore) *Service {
	return &Service{store: us}
}

// Register a new user.
func (s *Service) Register(mail, password string) error {
	if mail == "" || len(password) < 8 {
		return errors.New("invalid mail or password too short")
	}
	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}
	return s.store.Create(&models.User{
		Email:          mail,
		HashedPassword: hashed,
	})
}

// Login verifies credentials, returns a signed JWT.
func (s *Service) Login(mail, password string) (string, error) {
	u, err := s.store.FindByEmail(mail)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if !CheckPassword(u.HashedPassword, password) {
		return "", ErrInvalidCredentials
	}
	// CreateToken is from internal/auth/jwt.go
	return CreateToken(mail)
}

// Logout is now a no‑op: JWTs live client‑side until expiry.
func (s *Service) Logout(mail string) error {
	return nil
}

// Authorize parses and verifies the Bearer token.
func (s *Service) Authorize(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization") // “Bearer <token>”
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrUnauthorized
	}
	claims, err := ParseToken(parts[1])
	if err != nil {
		return "", ErrUnauthorized
	}
	return claims.Email, nil
}
