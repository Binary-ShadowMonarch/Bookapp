// internal/handlers/router.go

package handlers

import (
	"net/http"

	"bookapp/internal/auth" // import path
	"bookapp/internal/middleware"
)

func NewRouter(svc *auth.Service) http.Handler {
	mux := http.NewServeMux()

	// Public routes (with CORS)
	mux.Handle("/signup/request", middleware.CORS(RequestVerifyHandler(svc)))
	mux.Handle("/signup/verify", middleware.CORS(VerifyHandler(svc)))
	mux.Handle("/login", middleware.CORS(LoginHandler(svc)))
	mux.Handle("/refresh", middleware.CORS(RefreshHandler(svc))) // Add the new refresh route
	mux.Handle("/logout", middleware.CORS(LogoutHandler(svc)))

	// Protected routes (wrapped with JWT middleware)
	// The protected handler now represents any user-specific data endpoint, like /library.
	mux.Handle("/library", middleware.CORS(middleware.JWTMiddleware(svc, LibraryHandler(svc)))) // internal/handlers/router.go
	mux.Handle(
		"/upload",
		middleware.CORS(
			middleware.JWTMiddleware(svc,
				UploadHandler(svc),
			),
		),
	)
	return mux
}
