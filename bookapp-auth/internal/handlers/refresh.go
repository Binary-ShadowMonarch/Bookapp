// internal/handlers/refresh.go
package handlers

import (
	"bookapp/internal/auth"
	"net/http"
	"time"
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func RefreshHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get refresh token from cookie
		oldTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "unauthorized: no refresh token", http.StatusUnauthorized)
			return
		}

		// Rotate tokens
		newAccessToken, newRefreshToken, err := svc.Refresh(oldTokenCookie.Value)
		if err != nil {
			http.Error(w, "unauthorized: invalid refresh token", http.StatusUnauthorized)
			return
		}

		// Set new cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			Expires:  time.Now().UTC().Add(auth.AccessTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			Expires:  time.Now().UTC().Add(auth.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// Return JSON response for SvelteKit
		// response := RefreshResponse{
		// 	AccessToken:  newAccessToken,
		// 	RefreshToken: newRefreshToken,
		// 	ExpiresIn:    int64(auth.AccessTTL.Seconds()),
		// 	TokenType:    "Bearer",
		// }

		// w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		// json.NewEncoder(w).Encode(response)
	}
}
