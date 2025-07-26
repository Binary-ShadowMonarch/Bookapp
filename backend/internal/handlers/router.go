// internal/handlers/router.go
package handlers

import (
	"bookapp/internal/auth"
	"bookapp/internal/middleware"
	"net/http"
)

func NewRouter(svc *auth.Service) http.Handler {
	mux := http.NewServeMux()

	// Prefix all routes with /api
	api := http.NewServeMux()

	api.Handle("/register/request", middleware.CORS(RequestVerifyHandler(svc)))
	api.Handle("/register/verify", middleware.CORS(VerifyHandler(svc)))
	api.Handle("/login", middleware.CORS(LoginHandler(svc)))
	api.Handle("/refresh", middleware.CORS(RefreshHandler(svc)))
	api.Handle("/logout", middleware.CORS(LogoutHandler(svc)))
	api.Handle("/auth/google/login", middleware.CORS(GoogleLoginHandler(svc)))
	api.Handle("/auth/google/callback", middleware.CORS(GoogleCallbackHandler(svc)))

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.Handle("/library", LibraryHandler(svc))
	protectedMux.Handle("/library/", LibraryHandler(svc))

	// Upload routes
	protectedMux.Handle("/upload", UploadHandler(svc))

	// File proxy route - serves actual file content
	protectedMux.Handle("/files/", FileProxyHandler(svc))

	// User profile routes
	protectedMux.Handle("/profile", ProfileHandler(svc))

	api.Handle("/protected/",
		middleware.CORS(
			middleware.JWTMiddleware(svc,
				http.StripPrefix("/protected", protectedMux),
			),
		),
	)

	// Main router now uses the /api prefix
	mux.Handle("/api/", http.StripPrefix("/api", api))
	return mux
}
