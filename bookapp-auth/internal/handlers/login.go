// internal/handlers/login.go

package handlers

import (
	"bookapp/internal/auth"
	"fmt"
	"net/http"
	"time"
)

func LoginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		mail := r.FormValue("mail")
		pw := r.FormValue("password")

		accessToken, refreshToken, err := svc.Login(mail, pw)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set the long-lived refresh token as a secure, HttpOnly cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(auth.RefreshTTL),
			HttpOnly: true, // Prevents JavaScript access
			Secure:   true, // Sent only over HTTPS
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh", // Only sent to the /refresh endpoint
		})

		// Return the short-lived access token in the response body
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"accessToken":"%s"}`, accessToken)
	}
}
