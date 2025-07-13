package handlers

import (
	"fmt"
	"net/http"

	"bookapp/internal/auth"
)

func ProtectedHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := svc.Authorize(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "hello, %s", email)
	}
}
