package handlers

import (
	"bookapp/internal/auth"
	"log"
	"net/http"
	"time"
)

// RefreshResponse is what I send back when someone refreshes their tokens
// I'm not using this right now but I might need it later
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RefreshHandler handles token refresh requests
// when the access token expires, the frontend calls this to get new tokens
func RefreshHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Token refresh request received")
		
		// only allow POST requests for security
		if r.Method != http.MethodPost {
			log.Printf("DEBUG: Refresh method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get the refresh token from the cookie
		// this token is used to get new access and refresh tokens
		oldTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Printf("DEBUG: No refresh token found in cookie: %v", err)
			http.Error(w, "unauthorized: no refresh token", http.StatusUnauthorized)
			return
		}

		log.Printf("DEBUG: Attempting to refresh tokens")

		// exchange the old refresh token for new access and refresh tokens
		// this is called token rotation and it's more secure
		newAccessToken, newRefreshToken, err := svc.Refresh(oldTokenCookie.Value)
		if err != nil {
			log.Printf("DEBUG: Token refresh failed: %v", err)
			http.Error(w, "unauthorized: invalid refresh token", http.StatusUnauthorized)
			return
		}

		log.Printf("DEBUG: Tokens refreshed successfully, setting new cookies")

		// set the new access token as a secure cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			Expires:  time.Now().UTC().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// set the new refresh token as another secure cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			Expires:  time.Now().UTC().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// send back a success response (no content needed)
		// the new cookies are what matter
		log.Println("DEBUG: Token refresh complete")
		w.WriteHeader(http.StatusNoContent)
	}
}
