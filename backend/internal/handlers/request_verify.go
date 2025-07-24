package handlers

import (
	"bookapp/internal/auth"
	"fmt"
	"net/http"
)

func RequestVerifyHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		email := r.FormValue("mail")
		password := r.FormValue("password")
		if err := svc.RequestVerification(email, password); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Fprintln(w, "OK")
	}
}
