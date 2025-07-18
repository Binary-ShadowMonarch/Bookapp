// internal/handlers/login.go

package handlers

import (
	"bookapp/internal/auth"
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
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// after creating tokens…
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/", // valid on all API routes
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh", // only sent to your /refresh handler
		})
		// return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
