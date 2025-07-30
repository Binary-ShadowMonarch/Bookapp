package handlers

import (
	"bookapp/internal/auth"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// UploadResponse is what I send back after someone uploads a book
// right now it just has the URL, but I might add more info later
type UploadResponse struct {
	URL string `json:"url"`
	// you can add Size, Name, etc.
}

// UploadHandler handles book uploads
// this is how users add new books to their library
func UploadHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: Upload request received")
		
		// make sure the user is logged in before they can upload
		email, err := svc.Authorize(r)
		if err != nil {
			log.Printf("DEBUG: Upload unauthorized: %v", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// get the user info from my database
		u, err := svc.GetUserByEmail(email)
		if err != nil {
			log.Printf("DEBUG: User lookup failed for upload: %v", err)
			http.Error(w, "user lookup failed", http.StatusInternalServerError)
			return
		}

		log.Printf("DEBUG: User %d attempting to upload file", u.ID)

		// parse the multipart form to get the uploaded file
		// 20MB max file size limit
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			log.Printf("DEBUG: Failed to parse upload form: %v", err)
			http.Error(w, "could not parse form", http.StatusBadRequest)
			return
		}

		// get the actual file from the form
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("DEBUG: No file found in upload request: %v", err)
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("DEBUG: Uploading file %s (%d bytes) for user %d", header.Filename, header.Size, u.ID)

		// upload the file to my storage (MinIO)
		// this saves the book file in the user's personal folder
		url, err := svc.UploadFile(context.Background(), u.ID, file, header)
		if err != nil {
			log.Printf("DEBUG: Upload failed for user %d: %v", u.ID, err)
			http.Error(w, "upload failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("DEBUG: Successfully uploaded %s for user %d", header.Filename, u.ID)

		// send back the URL where the file was saved
		// the frontend might use this for something, I'm not sure yet
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			URL string `json:"url"`
		}{URL: url})
	}
}
