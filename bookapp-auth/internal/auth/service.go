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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"

	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/minio/minio-go/v7" // for minio.Client & IsBucketAlreadyOwnedByYou
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ErrInvalidCredentials is returned when login fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrUnauthorized is returned when a request lacks a valid JWT.
var ErrUnauthorized = errors.New("unauthorized")

var ErrEmailExistsLocal = errors.New("an account with this email already exists with a password, please log in with your password")
var ErrUserNotFound = errors.New("internal db server error")

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
	DeleteAllRefreshTokensForUser(userID int) error // ADD THIS LINE
}

// Service wraps auth logic.
type Service struct {
	store             UserStore
	minio             *minio.Client
	bucketPref        string
	googleOAuthConfig *oauth2.Config
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

	// 2) Construct Google OAuth2 Config
	googleOAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		// This is the URL Google will redirect to after the user signs in.
		// It must be one of the "Authorized redirect URIs" in your Google Cloud Console.
		RedirectURL: os.Getenv("GOOGLE_REDIRECT_URL"), // e.g., "http://localhost:8080/auth/google/callback"
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Service{
		store:             us,
		minio:             mc,
		bucketPref:        "user-",           // prefix for buckets
		googleOAuthConfig: googleOAuthConfig, // store it in the service
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
	expires := time.Now().UTC().Add(5 * time.Minute)

	// persist
	if err := s.store.SaveVerification(email, hashed, code, expires); err != nil {
		return err
	}

	// send email
	return s.sendVerificationEmail(email, code)
}

// func (s *Service) VerifyCode(email, code string) error {
// 	hashed, err := s.store.GetVerification(email, code)
// 	if err != nil {
// 		return ErrUnauthorized
// 	}

// 	if err := s.store.Create(&models.User{
// 		Email: email, HashedPassword: hashed, Provider: "local",
// 	}); err != nil {
// 		log.Printf("DEBUG: User creation failed: %v", err)
// 		return err
// 	}

// 	_ = s.store.DeleteVerification(email)

// 	u, _ := s.store.FindByEmail(email)
// 	bucket := s.bucketPref + strconv.Itoa(u.ID)

// 	ctx := context.Background()

// 	if ok, _ := s.minio.BucketExists(ctx, bucket); !ok {
// 		if err := s.minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
// 			return err
// 		}

// 		// 🔧 Set public read policy after creating the bucket
// 		policy := fmt.Sprintf(`{
// 			"Version": "2012-10-17",
// 			"Statement": [
// 				{
// 					"Action": ["s3:GetObject"],
// 					"Effect": "Allow",
// 					"Principal": "*",
// 					"Resource": ["arn:aws:s3:::%s/*"]
// 				}
// 			]
// 		}`, bucket)

// 		if err := s.minio.SetBucketPolicy(ctx, bucket, policy); err != nil {
// 			log.Printf("ERROR: Failed to set bucket policy: %v", err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// Login verifies credentials, returns a signed JWT.
// func (s *Service) Login(mail, password string) (accessToken, refreshToken string, err error) {
// 	log.Printf("DEBUG: LOGIN INIT ")
// 	u, err := s.store.FindByEmail(mail)
// 	if err != nil || !CheckPassword(u.HashedPassword, password) {
// 		return "", "", ErrInvalidCredentials
// 	}

// 	accessToken, err = CreateToken(u.Email, AccessTTL)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	log.Printf("accessToken created: %s", accessToken)
// 	refreshToken, err = CreateToken(u.Email, RefreshTTL)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Persist the new refresh token in the database
// 	expiresAt := time.Now().UTC().Add(RefreshTTL)
// 	log.Printf("refreshToken created: %s", refreshToken)
// 	if err := s.store.SaveRefreshToken(refreshToken, u.ID, expiresAt); err != nil {
// 		return "", "", err
// 	}

// 	return accessToken, refreshToken, nil
// }

func (s *Service) Refresh(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {

	log.Printf("DEBUG: Refresh executed")
	log.Printf("DEBUG: old Refresh token : %s", oldRefreshToken)

	// 1) Verify the JWT signature and expiration first
	if _, err := ParseToken(oldRefreshToken); err != nil {
		log.Printf("DEBUG: failed to parse token : %s", err)
		return "", "", ErrUnauthorized
	}

	log.Printf("DEBUG: token parsed")

	// a. Validate the old token in the DB
	userID, err := s.store.FindRefreshToken(oldRefreshToken)
	if err != nil {
		log.Printf("DEBUG: cant find the token : %s", err)
		return "", "", ErrUnauthorized
	}
	log.Printf("DEBUG: token validated")

	// c. Get user details to create new tokens
	u, err := s.store.FindByID(userID)
	if err != nil {
		log.Printf("DEBUG: critical cant find user( is the user logged in): %s", err)
		return "", "", err
	}

	// d. Issue a new pair of tokens
	newAccessToken, err = CreateToken(u.Email, AccessTTL)
	if err != nil {
		log.Printf("DEBUG: can't create the access token  : %s", err)
		return "", "", err
	}
	log.Printf("DEBUG: created access token")

	newRefreshToken, err = CreateToken(u.Email, RefreshTTL)
	if err != nil {
		log.Printf("DEBUG: can't create the refresh token  : %s", err)

		return "", "", err
	}
	log.Printf("DEBUG: created refresh token")

	// e. Persist the new refresh token
	expiresAt := time.Now().UTC().Add(RefreshTTL)
	if err := s.store.SaveRefreshToken(newRefreshToken, u.ID, expiresAt); err != nil {
		log.Printf("DEBUG: SaveRefreshToken failed  : %s", err)

		return "", "", err
	}
	log.Printf("DEBUG: tokens saved ")

	// b. Delete the old token (it has been used)
	if err := s.store.DeleteRefreshToken(oldRefreshToken); err != nil {
		log.Printf("DEBUG: cant delete refresh token : %s", err)
		return "", "", err
	}

	log.Printf("DEBUG:refresh token deleted  : %s", err)
	log.Printf("DEBUG: new accesstoken: %s \n new refreshtoken: %s", newAccessToken, newRefreshToken)

	return newAccessToken, newRefreshToken, nil
}

// Logout is now a no‑op: JWTs live client‑side until expiry.
func (s *Service) Logout(mail string) error {
	return nil
}

// Authorize parses and verifies the Bearer token or access_token cookie.
func (s *Service) Authorize(r *http.Request) (string, error) {
	var token string

	// 1. Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		token = parts[1]
	} else {
		// 2. Fallback to access_token cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			return "", ErrUnauthorized
		}
		token = cookie.Value
	}

	// 3. Parse and validate token
	claims, err := ParseToken(token)
	if err != nil {
		log.Printf("DEBUG : parse error :%s", err)
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

// UploadFile uploads the given multipart file to the MinIO bucket for the given user ID.
// It returns the public URL (or any URL scheme you choose) of the uploaded object.
func (s *Service) UploadFile(ctx context.Context, userID int, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Ensure bucket name matches your prefix + user ID
	bucket := s.bucketPref + strconv.Itoa(userID)

	// Create the bucket if it doesn't exist
	exists, err := s.minio.BucketExists(ctx, bucket)
	if err != nil {
		return "", fmt.Errorf("checking bucket existence: %w", err)
	}
	if !exists {
		if err := s.minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return "", fmt.Errorf("making bucket: %w", err)
		}
	}

	// Construct an object name: you might want to
	// prefix with a timestamp or user ID for uniqueness
	objectName := fmt.Sprintf("%d_%d%s",
		userID,
		time.Now().UTC().UnixNano(),
		filepath.Ext(header.Filename),
	)

	// Upload the file
	info, err := s.minio.PutObject(
		ctx,
		bucket,
		objectName,
		io.LimitReader(file, header.Size), // ensure PutObject knows size
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
			// You can set ACL-like via metadata or bucket policy
		},
	)
	if err != nil {
		return "", fmt.Errorf("uploading to minio: %w", err)
	}

	// Construct a URL – adjust to your MinIO endpoint / gateway
	url := fmt.Sprintf("http://%s/%s/%s", s.minio.EndpointURL().Host, bucket, info.Key)
	return url, nil
}

// GetUserByEmail is a thin wrapper around the underlying UserStore.
func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	return s.store.FindByEmail(email)
}

// ListFiles returns a slice of public URLs for every object
// in the given user's bucket.
func (s *Service) ListFiles(ctx context.Context, userID int) ([]string, error) {
	bucket := s.bucketPref + strconv.Itoa(userID)

	// make sure bucket exists
	exists, err := s.minio.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("checking bucket existence: %w", err)
	}
	if !exists {
		return nil, nil // no bucket → no files
	}
	// List all objects
	objectCh := s.minio.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	var urls []string
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, fmt.Errorf("listing objects: %w", obj.Err)
		}
		// build your public URL pattern
		urls = append(urls, fmt.Sprintf("http://%s/%s/%s",
			s.minio.EndpointURL().Host,
			bucket,
			obj.Key,
		))
	}
	return urls, nil
}

// DeleteFile removes a file from the user's bucket
func (s *Service) DeleteFile(ctx context.Context, userID int, fileName string) error {
	bucket := s.bucketPref + strconv.Itoa(userID)

	return s.minio.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{})
}

// GetFileInfo returns detailed information about a specific file
func (s *Service) GetFileInfo(ctx context.Context, userID int, fileName string) (*FileInfo, error) {
	bucket := s.bucketPref + strconv.Itoa(userID)

	objInfo, err := s.minio.StatObject(ctx, bucket, fileName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/%s/%s", s.minio.EndpointURL().Host, bucket, fileName)

	return &FileInfo{
		ID:       fileName,
		Name:     fileName,
		URL:      url,
		Size:     objInfo.Size,
		MimeType: objInfo.ContentType,
	}, nil
}

// FileInfo struct for file details
type FileInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Size     int64  `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// Enhanced Logout to revoke refresh tokens
func (s *Service) LogoutWithTokenRevocation(refreshToken string) error {
	if refreshToken != "" {
		// Revoke the specific refresh token
		return s.store.DeleteRefreshToken(refreshToken)
	}
	return nil
}

// LogoutAllSessions revokes all refresh tokens for a user
func (s *Service) LogoutAllSessions(email string) error {
	user, err := s.store.FindByEmail(email)
	if err != nil {
		return err
	}

	// TODO: Add method to store to delete all tokens for user
	return s.store.DeleteAllRefreshTokensForUser(user.ID)
}

// LoginOrRegisterWithGoogle handles the Google OAuth callback.
// It takes the authorization code from Google, validates it, and then either
// finds an existing Google user or creates a new one.
// It returns your app's own access and refresh tokens.
// LoginOrRegisterWithGoogle handles the Google OAuth callback.
// This is the updated and corrected version.
func (s *Service) LoginOrRegisterWithGoogle(ctx context.Context, code string) (accessToken, refreshToken string, err error) {
	// 1. Exchange code and validate Google's ID token (no changes here)
	googleToken, err := s.googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("failed to exchange code: %w", err)
	}
	rawIDToken, ok := googleToken.Extra("id_token").(string)
	if !ok {
		return "", "", errors.New("id_token not found in google token")
	}
	payload, err := idtoken.Validate(ctx, rawIDToken, s.googleOAuthConfig.ClientID)
	if err != nil {
		return "", "", fmt.Errorf("failed to validate id_token: %w", err)
	}
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return "", "", errors.New("email not found in claims")
	}

	// 2. Find or Create User in our Database
	user, err := s.store.FindByEmail(email)

	// Case A: User does NOT exist in our DB.
	if err != nil {
		// If the error is anything other than "not found", it's a real DB problem.
		if err != ErrUserNotFound {
			return "", "", err
		}

		// User is not found, so create them. This is their first sign-in.
		newUser := &models.User{
			Email:          email,
			Provider:       "google",
			ProviderID:     payload.Subject,
			HashedPassword: "", // No password needed for OAuth users
		}
		if err := s.store.Create(newUser); err != nil {
			return "", "", fmt.Errorf("failed to create google user: %w", err)
		}

		// Fetch the user again to get their database ID.
		user, err = s.store.FindByEmail(email)
		if err != nil {
			// This would be a critical error if it fails right after creation.
			return "", "", err
		}

		// Create their MinIO bucket since they are a new user.
		_ = s.createMinioBucketForUser(ctx, user.ID)

	} else {
		// Case B: User EXISTS in our DB.
		// Check if they signed up with a different provider (e.g., local password).
		if user.Provider != "google" {
			return "", "", ErrEmailExistsLocal
		}
		// If we're here, the user is a returning Google user.
		// The `user` variable is already populated and correct. We just proceed.
	}

	// 3. Issue Fresh Tokens
	// This step now runs for BOTH new and returning users, ensuring that
	// every successful Google login gets a new set of access/refresh tokens.
	return s.issueAndSaveTokens(user)
}

// Helper function to create MinIO bucket to avoid code duplication
func (s *Service) createMinioBucketForUser(ctx context.Context, userID int) error {
	bucket := s.bucketPref + strconv.Itoa(userID)
	if ok, _ := s.minio.BucketExists(ctx, bucket); !ok {
		if err := s.minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		policy := fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Action": ["s3:GetObject"],
                    "Effect": "Allow",
                    "Principal": "*",
                    "Resource": ["arn:aws:s3:::%s/*"]
                }
            ]
        }`, bucket)
		if err := s.minio.SetBucketPolicy(ctx, bucket, policy); err != nil {
			log.Printf("ERROR: Failed to set bucket policy: %v", err)
			return err
		}
	}
	return nil
}

// Helper function to issue and save tokens to avoid code duplication
func (s *Service) issueAndSaveTokens(u *models.User) (accessToken, refreshToken string, err error) {
	accessToken, err = CreateToken(u.Email, AccessTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = CreateToken(u.Email, RefreshTTL)
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().UTC().Add(RefreshTTL)
	if err := s.store.SaveRefreshToken(refreshToken, u.ID, expiresAt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Modify your Login function to use the new helper
func (s *Service) Login(mail, password string) (accessToken, refreshToken string, err error) {
	u, err := s.store.FindByEmail(mail)
	if err != nil || !CheckPassword(u.HashedPassword, password) {
		return "", "", ErrInvalidCredentials
	}

	if u.Provider != "local" {
		return "", "", fmt.Errorf("please sign in using %s", u.Provider)
	}

	return s.issueAndSaveTokens(u)
}

// And finally, modify the VerifyCode function to use the bucket creation helper
func (s *Service) VerifyCode(email, code string) error {
	hashed, err := s.store.GetVerification(email, code)
	if err != nil {
		return ErrUnauthorized
	}

	if err := s.store.Create(&models.User{
		Email: email, HashedPassword: hashed, Provider: "local",
	}); err != nil {
		log.Printf("DEBUG: User creation failed: %v", err)
		return err
	}

	_ = s.store.DeleteVerification(email)

	u, _ := s.store.FindByEmail(email)
	// Use the helper here
	return s.createMinioBucketForUser(context.Background(), u.ID)
}

// GoogleAuthCodeURL returns the URL for the Google consent page.
func (s *Service) GoogleAuthCodeURL(state string) string {
	return s.googleOAuthConfig.AuthCodeURL(state)
}
