// internal/auth/service.go
package auth

import (
	"bookapp/internal/models" // or "bookapp/internal/models" per your module
	// or "bookapp/internal/models" per your module
	"context"
	"errors" // for ErrInvalidCredentials,ErrUnauthorized, errors.New
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings" // for strings.SplitN
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

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
	// email verification management
	SaveVerification(email, hashedPw, code string, expiresAt time.Time) error
	GetVerification(email, code string) (hashedPw string, err error)
	DeleteVerification(email string) error
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

func (s *Service) RequestVerification(email, password string) error {
	if email == "" || len(password) < 8 {
		return errors.New("invalid")
	}

	// Check if user exists by email
	_, err := s.store.FindByEmail(email)
	if err != nil {
		// Check if it's a "user not found" error (you need to import your store package or use the specific error)
		// Replace this with your actual ErrUserNotFound from your store package
		if err.Error() == "user not found" || err.Error() == "sql: no rows in result set" {
			// User NOT found - continue with verification process (this is for new registrations)
		} else {
			// Some other database error
			return err
		}
	} else {
		// User EXISTS - return error to prevent duplicate registrations
		return errors.New("account associated with this email already exists")
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}

	// generate 6‑digit code
	code := fmt.Sprintf("%06d", rand.Intn(1_000_000))
	expires := time.Now().Add(15 * time.Minute)

	// persist
	if err := s.store.SaveVerification(email, hashed, code, expires); err != nil {
		return err
	}

	// send email
	return s.sendVerificationEmail(email, code)
}

func (s *Service) VerifyCode(email, code string) error {

	// log.Printf("DEBUG: VerifyCode called with email='%s', code='%s'", email, code)

	hashed, err := s.store.GetVerification(email, code)
	if err != nil {
		return ErrUnauthorized
	}

	// log.Printf("DEBUG: GetVerification succeeded, hashed password length: %d", len(hashed))

	// create real user
	if err := s.store.Create(&models.User{
		Email: email, HashedPassword: hashed, Provider: "local",
	}); err != nil {
		log.Printf("DEBUG: User creation failed: %v", err)
		return err
	}
	// log.Printf("DEBUG: User created successfully")
	// cleanup
	_ = s.store.DeleteVerification(email)
	// bucket
	u, _ := s.store.FindByEmail(email)
	bucket := s.bucketPref + strconv.Itoa(u.ID)
	ctx := context.Background()
	if ok, _ := s.minio.BucketExists(ctx, bucket); !ok {
		if err := s.minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
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

func (s *Service) sendVerificationEmail(to, code string) error {
	sg := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	from := mail.NewEmail("Books App", os.Getenv("SENDGRID_EMAIL"))
	subject := "Your verification code"
	toEmail := mail.NewEmail("", to)
	plainText := fmt.Sprintf("Your code is %s", code)
	htmlContent := fmt.Sprintf(`<p>Your verification code is <b>%s</b></p>`, code)
	message := mail.NewSingleEmail(from, subject, toEmail, plainText, htmlContent)
	_, err := sg.Send(message)
	// println(os.Getenv("SENDGRID_API_KEY"))
	return err
}
