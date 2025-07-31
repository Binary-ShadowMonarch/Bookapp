package handlers

import (
	"bookapp/internal/auth"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

// stateCookieName is the name of the cookie that stores the OAuth state
// this is for security to prevent CSRF attacks during Google login
const stateCookieName = "oauthstate"

// GoogleLoginHandler starts the Google OAuth flow
// this redirects users to Google's login page
func GoogleLoginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Google login initiated")

		// generate a random state string for security
		// this prevents CSRF attacks by making sure the callback comes from Google
		b := make([]byte, 16)
		rand.Read(b)
		state := base64.URLEncoding.EncodeToString(b)

		log.Printf("DEBUG: Generated OAuth state: %s", state)

		// store the state in a cookie so I can verify it later
		// this cookie expires in 10 minutes for security
		http.SetCookie(w, &http.Cookie{
			Name:     stateCookieName,
			Value:    state,
			Expires:  time.Now().UTC().Add(3 * time.Minute),
			HttpOnly: true,
			Path:     "/",
			Secure:   true, // needs HTTPS in production
			SameSite: http.SameSiteLaxMode,
		})

		// get the Google OAuth URL from my auth service
		// this URL will take users to Google's login page
		url := svc.GoogleAuthCodeURL(state)
		log.Printf("DEBUG: Redirecting to Google OAuth URL")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// GoogleCallbackHandler handles the redirect back from Google
// this is where Google sends users after they log in successfully
func GoogleCallbackHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Google OAuth callback received")

		// first, check the state cookie to make sure this is a legitimate callback
		// this prevents CSRF attacks
		stateCookie, err := r.Cookie(stateCookieName)
		if err != nil {
			log.Printf("DEBUG: OAuth state cookie not found: %v", err)
			http.Error(w, "state cookie not found", http.StatusBadRequest)
			return
		}

		// compare the state from the cookie with the state from the URL
		// they should match if this is a legitimate callback
		if r.URL.Query().Get("state") != stateCookie.Value {
			log.Printf("DEBUG: Invalid OAuth state - cookie: %s, URL: %s", stateCookie.Value, r.URL.Query().Get("state"))
			http.Error(w, "invalid oauth state", http.StatusBadRequest)
			return
		}

		log.Printf("DEBUG: OAuth state verified successfully")

		// get the authorization code from Google
		// this code is what I'll exchange for access tokens
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Println("DEBUG: No authorization code received from Google")
			http.Error(w, "no authorization code", http.StatusBadRequest)
			return
		}

		log.Printf("DEBUG: Exchanging authorization code for tokens")

		// exchange the code for access and refresh tokens
		// this also handles user registration if they're new
		accessToken, refreshToken, err := svc.LoginOrRegisterWithGoogle(r.Context(), code)
		if err != nil {
			// handle the special case where someone tries to use Google login
			// but they already have a local account with the same email
			if err == auth.ErrEmailExistsLocal {
				log.Printf("DEBUG: Google login failed - email already exists locally")
				http.Error(w, "A Local account exists with this email please login ", http.StatusBadRequest)

				// redirect to login page with an error message
				// http.Redirect(w, r, "/login?error=email_exists", http.StatusSeeOther)
				return
			}
			log.Printf("DEBUG: Google login failed: %v", err)
			http.Error(w, "failed to sign in with google: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("DEBUG: Google login successful, setting cookies")

		// set the access token as a secure cookie
		// this is the same as regular login
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().UTC().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// set the refresh token as another secure cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().UTC().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// redirect the user to the library page
		// they're now logged in and can access their books
		log.Printf("DEBUG: Redirecting user to library after successful Google login")
		http.Redirect(w, r,
			"/library",
			http.StatusSeeOther,
		)
	}
}
