package handlers

import (
	"bookapp/internal/auth"
	"fmt"
	"log"
	"net/http"
)

// VerifyHandler handles email verification
// this is used when someone registers and needs to verify their email
func VerifyHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Email verification request received")
		
		// only allow POST requests for security
		if r.Method != http.MethodPost {
			log.Printf("DEBUG: Verify method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// get the email and verification code from the form
		email := r.FormValue("mail")
		code := r.FormValue("code")
		
		log.Printf("DEBUG: Attempting to verify email %s with code", email)
		
		// verify the code with my auth service
		if err := svc.VerifyCode(email, code); err != nil {
			log.Printf("DEBUG: Email verification failed for %s: %v", email, err)
			http.Error(w, "invalid code", http.StatusUnauthorized)
			return
		}
		
		log.Printf("DEBUG: Email verification successful for %s", email)
		fmt.Fprintln(w, "verified")
	}
}
