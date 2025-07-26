// internal/handler/google.go
package handlers

import (
	"bookapp/internal/auth"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// stateCookieName is the name of the cookie that stores the OAuth state
const stateCookieName = "oauthstate"

// GoogleLoginHandler generates the Google OAuth URL and redirects the user.
func GoogleLoginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate a random state string for CSRF protection.
		b := make([]byte, 16)
		rand.Read(b)
		state := base64.URLEncoding.EncodeToString(b)

		// Store the state in a short-lived cookie.
		http.SetCookie(w, &http.Cookie{
			Name:     stateCookieName,
			Value:    state,
			Expires:  time.Now().Add(10 * time.Minute),
			HttpOnly: true,
			Path:     "/",
			Secure:   true, // Always true in production
			SameSite: http.SameSiteLaxMode,
		})

		// Get the authentication URL from our service config.
		// Note: We need to expose the config or a method to get the URL.
		// Let's create a helper in the service for this.
		url := svc.GoogleAuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// GoogleCallbackHandler handles the redirect from Google after user authentication.
func GoogleCallbackHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Check state cookie for CSRF.
		stateCookie, err := r.Cookie(stateCookieName)
		if err != nil {
			http.Error(w, "state cookie not found", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("state") != stateCookie.Value {
			http.Error(w, "invalid oauth state", http.StatusBadRequest)
			return
		}

		// 2. Get the authorization code from the query parameters.
		code := r.URL.Query().Get("code")

		// 3. Exchange code for tokens and handle user logic in the service.
		accessToken, refreshToken, err := svc.LoginOrRegisterWithGoogle(r.Context(), code)
		if err != nil {
			// Handle the specific error for existing local users
			if err == auth.ErrEmailExistsLocal {
				// Redirect to the login page with an error message
				// This is a user-friendly way to handle it.
				http.Redirect(w, r, "/login?error=email_exists", http.StatusSeeOther)
				return
			}
			http.Error(w, "failed to sign in with google: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 4. Set our application's tokens in secure cookies (same as local login).
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().UTC().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().UTC().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// 5. Redirect the user to the main application page.
		// The frontend will now have the tokens and can access protected routes.
		http.Redirect(w, r,
			"/library",
			http.StatusSeeOther,
		)
	}
}
