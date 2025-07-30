// internal/handlers/login.go
package handlers

import (
	"bookapp/internal/auth"
	"log"
	"net/http"
	"time"
)

// LoginHandler handles user login requests
// takes email and password from form data, validates them, and sets cookies
// this is probably the most important function since everyone needs to login
func LoginHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Login attempt received")
		
		// only allow POST requests for login (security thing)
		if r.Method != http.MethodPost {
			log.Println("DEBUG: Login failed - wrong HTTP method:", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get the email and password from the form data
		// these come from the frontend login form
		mail := r.FormValue("mail")
		pw := r.FormValue("password")
		
		log.Printf("DEBUG: Attempting login for email: %s", mail)

		// try to authenticate the user with my auth service
		// this will check if the email/password combo is correct
		accessToken, refreshToken, err := svc.Login(mail, pw)
		if err != nil {
			log.Printf("DEBUG: Login failed for %s: %v", mail, err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("DEBUG: Login successful for %s, setting cookies", mail)

		// set the access token as a secure cookie
		// this token is used to prove the user is logged in
		// HttpOnly means JavaScript can't steal it (security)
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().UTC().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   true, // needs HTTPS in production
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// set the refresh token as another secure cookie
		// this token is used to get a new access token when the old one expires
		// it lasts longer than the access token
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().UTC().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true, // needs HTTPS in production
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// send back a success response (no content needed)
		// the cookies are what matter, not the response body
		log.Println("DEBUG: Login complete, sending success response")
		w.WriteHeader(http.StatusNoContent)
	}
}
