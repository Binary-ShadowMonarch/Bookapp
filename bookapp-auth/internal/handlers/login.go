package handlers

import (
	"fmt"
	"net/http"

	"bookapp/internal/auth"
)

func LoginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		mail := r.FormValue("mail")
		pw := r.FormValue("password")
		token, err := svc.Login(mail, pw)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// return JSON instead of cookies:
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"token":"%s"}`, token)
	}
}
