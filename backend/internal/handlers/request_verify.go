package handlers

import (
	"bookapp/internal/auth"
	"fmt"
	"log"
	"net/http"
)

// RequestVerifyHandler handles email verification requests
// this sends a verification code to someone's email when they register
func RequestVerifyHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Email verification request initiated")
		
		// only allow POST requests for security
		if r.Method != http.MethodPost {
			log.Printf("DEBUG: Request verify method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// get the email and password from the form
		email := r.FormValue("mail")
		password := r.FormValue("password")
		
		log.Printf("DEBUG: Requesting verification email for %s", email)
		
		// send the verification email through my auth service
		if err := svc.RequestVerification(email, password); err != nil {
			log.Printf("DEBUG: Verification request failed for %s: %v", email, err)
			http.Error(w, err.Error(), 400)
			return
		}
		
		log.Printf("DEBUG: Verification email sent successfully to %s", email)
		fmt.Fprintln(w, "OK")
	}
}
