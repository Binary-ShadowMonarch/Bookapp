// internal/handlers/router.go

package handlers

import (
	"net/http"

	"bookapp/internal/auth" // import path
	"bookapp/internal/middleware"
)

func NewRouter(svc *auth.Service) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.Handle("/register", RegisterHandler(svc))
	mux.Handle("/login", LoginHandler(svc))
	mux.Handle("/refresh", RefreshHandler(svc)) // Add the new refresh route
	mux.Handle("/logout", LogoutHandler(svc))

	// Protected routes (wrapped with JWT middleware)
	// The protected handler now represents any user-specific data endpoint, like /library.
	mux.Handle("/protected", middleware.JWTMiddleware(svc, ProtectedHandler(svc)))

	return mux
}
