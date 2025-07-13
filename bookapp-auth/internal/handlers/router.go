package handlers

import (
	"net/http"

	"bookapp/internal/auth"
)

func NewRouter(svc *auth.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/register", RegisterHandler(svc))
	mux.Handle("/login", LoginHandler(svc))
	mux.Handle("/logout", LogoutHandler(svc))
	mux.Handle("/protected", ProtectedHandler(svc))
	return mux
}
