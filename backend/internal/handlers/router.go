// internal/handlers/router.go
package handlers

import (
	"bookapp/internal/auth"
	"bookapp/internal/middleware"
	"log"
	"net/http"
)

// NewRouter creates the main HTTP router for this book app
// this is where I wire up all the endpoints and middleware
// think of it as the traffic controller for all incoming requests
func NewRouter(svc *auth.Service) http.Handler {
	log.Println("DEBUG: Setting up router - this is where all the magic happens")

	mux := http.NewServeMux()

	// all my API routes go under /api prefix to keep things organized
	// like having separate folders for different types of stuff
	api := http.NewServeMux()

	// these are the public routes anyone can access (no login needed)
	// registration and login stuff basically
	log.Println("DEBUG: Setting up public routes (registration, login, etc.)")
	api.Handle("/healthz", middleware.CORS(HealthHandler()))
	api.Handle("/register/request", middleware.CORS(RequestVerifyHandler(svc)))
	api.Handle("/register/verify", middleware.CORS(VerifyHandler(svc)))
	api.Handle("/login", middleware.CORS(LoginHandler(svc)))
	api.Handle("/refresh", middleware.CORS(RefreshHandler(svc)))
	api.Handle("/logout", middleware.CORS(LogoutHandler(svc)))
	api.Handle("/auth/google/login", middleware.CORS(GoogleLoginHandler(svc)))
	api.Handle("/auth/google/callback", middleware.CORS(GoogleCallbackHandler(svc)))

	// these routes need authentication - user must be logged in
	// I'll wrap them with JWT middleware to check if they have valid tokens
	log.Println("DEBUG: Setting up protected routes (need login)")
	protectedMux := http.NewServeMux()
	protectedMux.Handle("/library", LibraryHandler(svc))
	protectedMux.Handle("/library/", LibraryHandler(svc))

	// upload route for adding new books to the library
	protectedMux.Handle("/upload", UploadHandler(svc))

	// this serves the actual book files (epub, pdf, etc.)
	// users can only access files they own
	protectedMux.Handle("/files/", FileProxyHandler(svc))

	// user profile stuff - view/edit account info
	protectedMux.Handle("/profile", ProfileHandler(svc))

	// wrap all protected routes with JWT middleware
	// this checks if the user has a valid token before letting them access anything
	api.Handle("/protected/",
		middleware.CORS(
			middleware.JWTMiddleware(svc,
				http.StripPrefix("/protected", protectedMux),
			),
		),
	)

	// finally, mount everything under /api prefix
	// so all my routes will be like /api/login, /api/protected/library, etc.
	log.Println("DEBUG: Router setup complete - all routes mounted under /api")
	mux.Handle("/api/", http.StripPrefix("/api", api))
	return mux
}
