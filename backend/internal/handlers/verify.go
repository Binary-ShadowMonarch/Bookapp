package handlers

import (
	"bookapp/internal/auth"
	"fmt"
	"net/http"
)

func VerifyHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		email := r.FormValue("mail")
		code := r.FormValue("code")
		if err := svc.VerifyCode(email, code); err != nil {
			http.Error(w, "invalid code", http.StatusUnauthorized)
			return
		}
		fmt.Fprintln(w, "verified")
	}
}
