// internal/handlers/router.go
package handlers

import (
	"bookapp/internal/auth"
	"bookapp/internal/middleware"
	"net/http"
)

func NewRouter(svc *auth.Service) http.Handler {
	mux := http.NewServeMux()

	// Public routes (with CORS)
	mux.Handle("/signup/request", middleware.CORS(RequestVerifyHandler(svc)))
	mux.Handle("/signup/verify", middleware.CORS(VerifyHandler(svc)))
	mux.Handle("/login", middleware.CORS(LoginHandler(svc)))
	mux.Handle("/refresh", middleware.CORS(RefreshHandler(svc)))
	mux.Handle("/logout", middleware.CORS(LogoutHandler(svc)))

	// This route starts the Google sign-in process.
	mux.Handle("/auth/google/login", middleware.CORS(GoogleLoginHandler(svc)))
	// This is the callback URL that you configured in your Google Cloud Console.
	mux.Handle("/auth/google/callback", middleware.CORS(GoogleCallbackHandler(svc)))

	// Protected routes under /protected prefix
	protectedMux := http.NewServeMux()

	// Library routes
	protectedMux.Handle("/library", LibraryHandler(svc))
	protectedMux.Handle("/library/", LibraryHandler(svc)) // handles /library/* paths

	// Upload routes
	protectedMux.Handle("/upload", UploadHandler(svc))

	// File management routes
	protectedMux.Handle("/files", ListFilesHandler(svc))
	protectedMux.Handle("/files/", FileHandler(svc)) // handles /files/{id} operations

	// User profile routes
	protectedMux.Handle("/profile", ProfileHandler(svc))

	// Wrap all protected routes with middleware
	mux.Handle("/protected/",
		middleware.CORS(
			middleware.JWTMiddleware(svc,
				http.StripPrefix("/protected", protectedMux),
			),
		),
	)

	return mux
}
