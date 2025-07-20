// internal/handlers/logout.go
package handlers

import (
	"bookapp/internal/auth"
	"net/http"
	"time"
)

func LogoutHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get refresh token to revoke it
		refreshCookie, err := r.Cookie("refresh_token")
		if err == nil {
			// Revoke the refresh token from database
			_ = svc.LogoutWithTokenRevocation(refreshCookie.Value)
		}

		// Clear cookies by setting them to expire immediately
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		w.WriteHeader(http.StatusNoContent)
	}
}
