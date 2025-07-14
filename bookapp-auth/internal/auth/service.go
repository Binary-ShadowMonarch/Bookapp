// internal/auth/service.go
package auth

import (
	"bookapp/internal/models" // or "bookapp/internal/models" per your module
	"context"
	"errors" // for ErrInvalidCredentials,ErrUnauthorized, errors.New
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings" // for strings.SplitN
	"time"

	"github.com/minio/minio-go/v7" // for minio.Client & IsBucketAlreadyOwnedByYou
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ErrInvalidCredentials is returned when login fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrUnauthorized is returned when a request lacks a valid JWT.
var ErrUnauthorized = errors.New("unauthorized")

// UserStore defines the methods we need for your store.
type UserStore interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindByID(int) (*models.User, error) // Add this
	Update(*models.User) error

	// Add these for refresh tokens
	SaveRefreshToken(token string, userID int, expiresAt time.Time) error
	DeleteRefreshToken(token string) error
	FindRefreshToken(token string) (int, error)
}

// Service wraps auth logic.
type Service struct {
	store      UserStore
	minio      *minio.Client
	bucketPref string
}

// NewService constructs the auth Service.
func NewService(us UserStore) *Service {
	// 1) Construct MinIO client
	endpoint := os.Getenv("MINIO_ENDPOINT") // e.g. "play.min.io:9000"
	accessKey := os.Getenv("MINIO_KEY")
	secretKey := os.Getenv("MINIO_SECRET")

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("unable to initialize MinIO client: %v", err)
	}

	return &Service{
		store:      us,
		minio:      mc,
		bucketPref: "user-", // prefix for buckets
	}
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
	if err := s.store.Create(&models.User{
		Email:          mail,
		HashedPassword: hashed,
		Provider:       "local",
	}); err != nil {
		return err
	}

	// --- NEW: fetch the user to get its ID ---

	u, err := s.store.FindByEmail(mail)
	if err != nil {
		return err
	}

	bucketName := s.bucketPref + strconv.Itoa(u.ID)
	ctx := context.Background()

	// Check existence
	exists, err := s.minio.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket %q exists: %w", bucketName, err)
	}
	// Create if missing
	if !exists {
		if err := s.minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("could not create bucket %q: %w", bucketName, err)
		}
	}
	print("sucessfully registered")
	return nil

}

// Login verifies credentials, returns a signed JWT.
func (s *Service) Login(mail, password string) (accessToken, refreshToken string, err error) {
	u, err := s.store.FindByEmail(mail)
	if err != nil || !CheckPassword(u.HashedPassword, password) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err = CreateToken(u.Email, AccessTTL)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = CreateToken(u.Email, RefreshTTL)
	if err != nil {
		return "", "", err
	}

	// Persist the new refresh token in the database
	expiresAt := time.Now().Add(RefreshTTL)
	if err := s.store.SaveRefreshToken(refreshToken, u.ID, expiresAt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// a. Validate the old token in the DB
	userID, err := s.store.FindRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", ErrUnauthorized
	}

	// b. Delete the old token (it has been used)
	if err := s.store.DeleteRefreshToken(oldRefreshToken); err != nil {
		return "", "", err
	}

	// c. Get user details to create new tokens
	u, err := s.store.FindByID(userID)
	if err != nil {
		return "", "", err
	}

	// d. Issue a new pair of tokens
	newAccessToken, err = CreateToken(u.Email, AccessTTL)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err = CreateToken(u.Email, RefreshTTL)
	if err != nil {
		return "", "", err
	}

	// e. Persist the new refresh token
	expiresAt := time.Now().Add(RefreshTTL)
	if err := s.store.SaveRefreshToken(newRefreshToken, u.ID, expiresAt); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
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
