// /internal/handlers/
package handlers

import (
	"fmt"
	"net/http"

	"bookapp/internal/auth"
)

func RegisterHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		mail := r.FormValue("mail")
		pw := r.FormValue("password")
		if err := svc.Register(mail, pw); err != nil {
			if err.Error() == "user already exists" {
				http.Error(w, "user already exists", http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		fmt.Fprintln(w, "registered")
	}
}
