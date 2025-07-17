// internal/handlers/refresh.go

package handlers

import (
	"net/http"
	"time"

	"bookapp/internal/auth" // Adjust import path
)

func RefreshHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the refresh token from the cookie
		oldTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "unauthorized: no refresh token", http.StatusUnauthorized)
			return
		}

		// 2. Call the service to rotate the tokens
		newAccessToken, newRefreshToken, err := svc.Refresh(oldTokenCookie.Value)
		if err != nil {
			http.Error(w, "unauthorized: no refresh token", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			Expires:  time.Now().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/", // valid on all API routes
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			Expires:  time.Now().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh", // only sent to your /refresh handler
		})
		// return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
