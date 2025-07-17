package handlers

import (
	"net/http"

	"bookapp/internal/auth"
)

func LogoutHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Optional: validate token if you want extra safety
		if _, err := svc.Authorize(r); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Invalidate cookies on the client by setting them with past expiry
		clearCookie := func(name string, path string) {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     path,
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}

		clearCookie("access_token", "/")
		clearCookie("refresh_token", "/refresh")

		// Also remove server-side refresh token (optional, already done)
		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}
