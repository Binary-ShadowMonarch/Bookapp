// internal/handlers/logout.go
package handlers

import (
	"bookapp/internal/auth"
	"log"
	"net/http"
	"time"
)

// LogoutHandler handles user logout requests
// this clears the user's session and invalidates their tokens
func LogoutHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Logout request received")

		// only allow POST requests for security
		if r.Method != http.MethodPost {
			log.Printf("DEBUG: Logout method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// try to get the refresh token from the cookie
		// if it exists, I'll revoke it in my database
		refreshCookie, err := r.Cookie("refresh_token")
		if err == nil {
			log.Printf("DEBUG: Revoking refresh token")
			// revoke the refresh token so it can't be used again
			_ = svc.LogoutWithTokenRevocation(refreshCookie.Value)
		} else {
			log.Printf("DEBUG: No refresh token found to revoke: %v", err)
		}

		log.Printf("DEBUG: Clearing authentication cookies")

		// clear the access token cookie by setting it to expire in the past
		// this effectively deletes the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().UTC().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   false, // false for logout since we want it to work without HTTPS
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// clear the refresh token cookie the same way
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   false, // false for logout since we want it to work without HTTPS
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		log.Println("DEBUG: Logout complete")
		w.WriteHeader(http.StatusNoContent)
	}
}
