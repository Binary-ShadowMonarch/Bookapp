package handlers

import (
	"bookapp/internal/auth"
	"context"
	"encoding/json"
	"net/http"
)

func LibraryHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) authorize
		email, err := svc.Authorize(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// 2) lookup user to get ID
		u, err := svc.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "user lookup failed", http.StatusInternalServerError)
			return
		}
		// 3) list files
		urls, err := svc.ListFiles(context.Background(), u.ID)
		if err != nil {
			http.Error(w, "could not list files: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// 4) respond
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urls)
	}
}
