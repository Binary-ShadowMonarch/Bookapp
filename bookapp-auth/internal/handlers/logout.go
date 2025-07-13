package handlers

import (
	"fmt"
	"net/http"

	"bookapp/internal/auth"
)

func LogoutHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Optional: ensure token is valid so clients can’t “logout” random users
		if _, err := svc.Authorize(r); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Nothing to clear server‑side; just respond success.
		// Front‑end should drop its copy of the JWT (e.g. clear the store).
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintln(w, "logged out")
	}
}
