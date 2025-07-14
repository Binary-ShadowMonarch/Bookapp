// internal/handlers/refresh.go

package handlers

import (
	"fmt"
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
			// If refresh fails, clear the cookie to force re-login
			http.SetCookie(w, &http.Cookie{
				Name:   "refresh_token",
				Value:  "",
				MaxAge: -1,
			})
			http.Error(w, "unauthorized: invalid refresh token", http.StatusUnauthorized)
			return
		}

		// 3. Set the new rotated refresh token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			Expires:  time.Now().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh",
		})

		// 4. Return the new access token in the response body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"accessToken":"%s"}`, newAccessToken)
	}
}
