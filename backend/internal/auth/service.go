package auth

import (
	"bookapp/internal/models" // or "bookapp/internal/models" per your module
	"bookapp/internal/store"

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

// these are the error messages I return when something goes wrong
// I use these to tell the frontend what happened
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUnauthorized = errors.New("unauthorized")
var ErrEmailExistsLocal = errors.New("an account with this email already exists, please log in with your password")
var ErrUserNotFound = errors.New("user not found")

// UserStore defines all the database operations I need
// this is an interface so I can easily swap out different databases later
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
	DeleteAllRefreshTokensForUser(userID int) error
}

// Service contains all my authentication and file storage logic
// this is the main service that handles login, registration, file uploads, etc.
type Service struct {
	store             UserStore
	minio             *minio.Client
	bucketPref        string
	googleOAuthConfig *oauth2.Config
}

// NewService creates a new auth service with all the necessary connections
// this sets up MinIO for file storage and Google OAuth for login
func NewService(us UserStore) *Service {
	log.Println("DEBUG: Initializing auth service")

	// set up MinIO client for file storage
	// MinIO is like AWS S3 but I can run it locally
	endpoint := os.Getenv("MINIO_ENDPOINT") // e.g. "play.min.io:9000"
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	log.Printf("DEBUG: Connecting to MinIO at %s", endpoint)
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("CRITICAL: Failed to initialize MinIO client: %v", err)
	}

	// set up Google OAuth configuration
	// this is what allows users to login with their Google account
	log.Println("DEBUG: Setting up Google OAuth configuration")
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

	log.Println("DEBUG: Auth service initialization complete")
	return &Service{
		store:             us,
		minio:             mc,
		bucketPref:        "user-",           // prefix for buckets
		googleOAuthConfig: googleOAuthConfig, // store it in the service
	}
}

// RequestVerification handles the first step of user registration
// this sends a verification email with a code to confirm the user's email address
func (s *Service) RequestVerification(email, password string) error {
	log.Printf("DEBUG: Requesting verification for email: %s", email)

	// validate the input - email can't be empty, password must be at least 8 characters
	if email == "" || len(password) < 8 {
		log.Printf("DEBUG: Invalid credentials for verification - email: %s, password length: %d", email, len(password))
		return errors.New("invalid credentials")
	}

	// check if a user with this email already exists
	u, err := s.store.FindByEmail(email)
	if err != nil {
		// if the error is "user not found", that's good - it means this is a new registration
		if err == store.ErrUserNotFound {
			log.Printf("DEBUG: New user registration for %s", email)
			// user NOT found - continue with verification process (this is for new registrations)
		} else {
			// some other database error occurred
			log.Printf("DEBUG: Database error during verification request: %v", err)
			return err
		}
	} else {
		log.Printf("DEBUG: User already exists for %s with provider %s", email, u.Provider)
		if u.Provider == "local" {

			return fmt.Errorf("ACCOUNT ASSOCIATED WITH THIS EMAIL EXISTS SIGN IN USING PASSWORD")
		}
		return fmt.Errorf("ACCOUNT ASSOCIATED WITH THIS EMAIL EXISTS SIGN IN USING %s", strings.ToUpper(u.Provider))
		// user EXISTS - return error to prevent duplicate registrations
	}

	// hash the password before storing it
	hashed, err := HashPassword(password)
	if err != nil {
		log.Printf("DEBUG: Failed to hash password for %s: %v", email, err)
		return err
	}

	// generate a random 6-digit verification code
	code := fmt.Sprintf("%06d", rand.Intn(1_000_000))
	expires := time.Now().UTC().Add(5 * time.Minute) // code expires in 5 minutes

	log.Printf("DEBUG: Generated verification code for %s, expires at %v", email, expires)

	// save the verification data to the database
	if err := s.store.SaveVerification(email, hashed, code, expires); err != nil {
		log.Printf("DEBUG: Failed to save verification for %s: %v", email, err)
		return err
	}

	// send the verification email to the user
	log.Printf("DEBUG: Sending verification email to %s", email)
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

// Refresh exchanges an old refresh token for new access and refresh tokens
// this is called token rotation and it's more secure than just extending tokens
func (s *Service) Refresh(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	log.Printf("DEBUG: Token refresh initiated")
	log.Printf("DEBUG: Old refresh token: %s", oldRefreshToken)

	// first, verify the JWT signature and check if it's expired
	if _, err := ParseToken(oldRefreshToken); err != nil {
		log.Printf("DEBUG: Failed to parse refresh token: %v", err)
		return "", "", ErrUnauthorized
	}

	log.Printf("DEBUG: Refresh token parsed successfully")

	// check if this refresh token exists in my database
	// this prevents using old tokens that have already been used
	userID, err := s.store.FindRefreshToken(oldRefreshToken)
	if err != nil {
		log.Printf("DEBUG: Refresh token not found in database: %v", err)
		return "", "", ErrUnauthorized
	}
	log.Printf("DEBUG: Refresh token validated in database")

	// get the user details so I can create new tokens for them
	u, err := s.store.FindByID(userID)
	if err != nil {
		log.Printf("DEBUG: Critical error - can't find user for refresh: %v", err)
		return "", "", err
	}

	// create a new access token (short-lived)
	newAccessToken, err = CreateToken(u.Email, AccessTTL)
	if err != nil {
		log.Printf("DEBUG: Failed to create new access token: %v", err)
		return "", "", err
	}
	log.Printf("DEBUG: New access token created successfully")

	// create a new refresh token (long-lived)
	newRefreshToken, err = CreateToken(u.Email, RefreshTTL)
	if err != nil {
		log.Printf("DEBUG: Failed to create new refresh token: %v", err)
		return "", "", err
	}
	log.Printf("DEBUG: New refresh token created successfully")

	// save the new refresh token to the database
	expiresAt := time.Now().UTC().Add(RefreshTTL)
	if err := s.store.SaveRefreshToken(newRefreshToken, u.ID, expiresAt); err != nil {
		log.Printf("DEBUG: Failed to save new refresh token: %v", err)
		return "", "", err
	}
	log.Printf("DEBUG: New refresh token saved to database")

	// delete the old refresh token since it's been used
	// this prevents the same token from being used multiple times
	if err := s.store.DeleteRefreshToken(oldRefreshToken); err != nil {
		log.Printf("DEBUG: Failed to delete old refresh token: %v", err)
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
// Authorize extracts and validates the JWT token from a request
// this is used by middleware to check if a user is logged in
// it looks for the token in either the Authorization header or a cookie(for sveltekit integration)
func (s *Service) Authorize(r *http.Request) (string, error) {
	log.Printf("DEBUG: Authorizing request: %s %s", r.Method, r.URL.Path)
	var token string

	// first, try to get the token from the Authorization header
	// found this is the standard way APIs handle authentication
	authHeader := r.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		token = parts[1]
		log.Printf("DEBUG: Found token in Authorization header")
	} else {
		// if no Authorization header, try to get the token from a cookie
		// this is useful for web applications(browser cookie)
		cookie, err := r.Cookie("access_token")
		if err != nil {
			log.Printf("DEBUG: No token found in header or cookie")
			return "", ErrUnauthorized
		}
		token = cookie.Value
		log.Printf("DEBUG: Found token in cookie")
	}

	// parse and validate the token
	claims, err := ParseToken(token)
	if err != nil {
		log.Printf("DEBUG: Token parsing failed: %v", err)
		return "", ErrUnauthorized
	}

	log.Printf("DEBUG: Authorization successful for %s", claims.Email)
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

// UploadFile uploads a book file to the user's personal storage bucket
// this is how users add new books to their library
func (s *Service) UploadFile(ctx context.Context, userID int, file multipart.File, header *multipart.FileHeader) (string, error) {
	log.Printf("DEBUG: Uploading file %s for user %d", header.Filename, userID)

	// create the bucket name for this user
	// each user gets their own bucket to keep their files separate
	bucket := s.bucketPref + strconv.Itoa(userID)

	// make sure the user's bucket exists, create it if it doesn't
	exists, err := s.minio.BucketExists(ctx, bucket)
	if err != nil {
		log.Printf("DEBUG: Failed to check bucket existence: %v", err)
		return "", fmt.Errorf("checking bucket existence: %w", err)
	}
	if !exists {
		log.Printf("DEBUG: Creating bucket for user %d", userID)
		if err := s.minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			log.Printf("DEBUG: Failed to create bucket: %v", err)
			return "", fmt.Errorf("making bucket: %w", err)
		}
	}

	// create a unique filename to avoid conflicts
	// I use userID + timestamp + original extension
	objectName := fmt.Sprintf("%d_%d%s",
		userID,
		time.Now().UTC().UnixNano(),
		filepath.Ext(header.Filename),
	)

	log.Printf("DEBUG: Uploading %s as %s to bucket %s", header.Filename, objectName, bucket)

	// actually upload the file to MinIO
	info, err := s.minio.PutObject(
		ctx,
		bucket,
		objectName,
		io.LimitReader(file, header.Size), // make sure MinIO knows the file size
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
			// I can set ACL-like permissions via metadata or bucket policy
		},
	)
	if err != nil {
		log.Printf("DEBUG: Failed to upload file to MinIO: %v", err)
		return "", fmt.Errorf("uploading to minio: %w", err)
	}

	// create a URL that goes through my server instead of direct MinIO access
	// this gives me control over who can access the files
	url := fmt.Sprintf("/api/protected/files/%s/%s", bucket, info.Key)
	log.Printf("DEBUG: File uploaded successfully, URL: %s", url)
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
		log.Printf("DEBUG: FindByEmail error: %v (type: %T)", err, err)
		if err != store.ErrUserNotFound {
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

// Login handles user authentication with email and password
// this is the main login function that validates credentials and issues tokens
func (s *Service) Login(mail, password string) (accessToken, refreshToken string, err error) {
	log.Printf("DEBUG: Login attempt for email: %s", mail)

	// first, try to find the user by email
	u, err := s.store.FindByEmail(mail)

	// check if the user exists - if not, login fails
	if err != nil {
		if err == store.ErrUserNotFound {
			log.Printf("DEBUG: Login failed - user not found: %s", mail)
			return "", "", ErrUserNotFound
		}
		log.Printf("DEBUG: Login failed - database error: %v", err)
		return "", "", ErrInvalidCredentials
	}

	// check if this is a local account (not Google login)
	// if someone tries to login with password but they used Google, tell them to use Google
	if u.Provider != "local" {
		log.Printf("DEBUG: Login failed - account uses %s provider, not local", u.Provider)
		return "", "", fmt.Errorf("ACCOUNT ASSOCIATED WITH THIS EMAIL EXISTS SIGN IN USING %s", strings.ToUpper(u.Provider))
	}

	// verify the password is correct
	if !CheckPassword(u.HashedPassword, password) {
		log.Printf("DEBUG: Login failed - invalid password for %s", mail)
		return "", "", ErrInvalidCredentials
	}

	log.Printf("DEBUG: Login successful for %s, issuing tokens", mail)
	// if all checks pass, create and return the tokens
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

func (s *Service) GetFileStream(ctx context.Context, userID int, fileName string) (*minio.Object, error) {
	bucket := s.bucketPref + strconv.Itoa(userID)
	return s.minio.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
}
