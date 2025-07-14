package middleware

import (
	"net/http"

	"bookapp/internal/auth"
)

// JWTMiddleware ensures that a valid Bearer token is present.
// If valid, it calls next; otherwise returns 401.
func JWTMiddleware(svc *auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// svc.Authorize returns the email from the token or error
		if _, err := svc.Authorize(r); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
