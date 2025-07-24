// internal/handlers/upload.go
package handlers

import (
	"bookapp/internal/auth"
	"context"
	"encoding/json"
	"net/http"
)

// UploadResponse is the JSON response returned after uploading.
type UploadResponse struct {
	URL string `json:"url"`
	// you can add Size, Name, etc.
}

// UploadHandler handles POST /upload requests.
func UploadHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Authorize the user
		email, err := svc.Authorize(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Lookup user via the new helper
		u, err := svc.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "user lookup failed", http.StatusInternalServerError)
			return
		}

		// 3. Parse form, pull out file, etc…
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			http.Error(w, "could not parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 4. Upload to MinIO under u.ID
		url, err := svc.UploadFile(context.Background(), u.ID, file, header)
		if err != nil {
			http.Error(w, "upload failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Return the URL
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			URL string `json:"url"`
		}{URL: url})
	}
}
