// internal/handlers/library.go
package handlers

import (
	"bookapp/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// LibraryResponse is what I send back when someone asks for their library
// contains a list of files and the total count
type LibraryResponse struct {
	Files []FileInfo `json:"files"`
	Total int        `json:"total"`
}

// FileInfo represents a single book file in the library
// this is what gets shown in the frontend library view
type FileInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Size     int64  `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// LibraryHandler handles all library-related requests
// right now it only supports GET to list all books, but I might add more later
func LibraryHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Library request received: %s %s", r.Method, r.URL.Path)
		
		switch r.Method {
		case http.MethodGet:
			handleGetLibrary(svc, w, r)
		default:
			log.Printf("DEBUG: Library method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleGetLibrary gets all the books for a logged-in user
// this is the main function that shows someone their personal library
func handleGetLibrary(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Getting library for user")
	
	// get the user from the JWT token (set by middleware)
	email, err := svc.Authorize(r)
	if err != nil {
		log.Printf("DEBUG: Library access unauthorized: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// look up the user in my database
	user, err := svc.GetUserByEmail(email)
	if err != nil {
		log.Printf("DEBUG: User lookup failed for %s: %v", email, err)
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Getting files for user ID: %d", user.ID)

	// get all the file URLs from my storage (MinIO)
	urls, err := svc.ListFiles(context.Background(), user.ID)
	if err != nil {
		log.Printf("DEBUG: Failed to list files for user %d: %v", user.ID, err)
		http.Error(w, "could not list files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Found %d files for user %d", len(urls), user.ID)

	// convert the raw URLs into nice FileInfo objects
	// I also create proxy URLs so users can't see my MinIO setup
	files := make([]FileInfo, len(urls))
	for i, url := range urls {
		// extract the filename from the MinIO URL
		parts := strings.Split(url, "/")
		filename := "unknown"
		if len(parts) > 0 {
			filename = parts[len(parts)-1]
		}

		// create a proxy URL that goes through my server instead of direct MinIO access
		proxyURL := fmt.Sprintf("/api/protected/files/user-%d/%s", user.ID, filename)

		files[i] = FileInfo{
			ID:   strconv.Itoa(i), // using index as ID for now, might change later
			Name: filename,
			URL:  proxyURL, // this is the URL the frontend will use to read the book
		}
	}

	// send back the library data as JSON
	response := LibraryResponse{
		Files: files,
		Total: len(files),
	}

	log.Printf("DEBUG: Sending library response with %d files", len(files))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// FileHandler handles individual file operations like delete and get details
// this is for when someone wants to do something with a specific book
func FileHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: File operation request: %s %s", r.Method, r.URL.Path)
		
		switch r.Method {
		case http.MethodDelete:
			handleDeleteFile(svc, w, r)
		case http.MethodGet:
			handleGetFile(svc, w, r)
		default:
			log.Printf("DEBUG: File method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleDeleteFile deletes a book from someone's library
// they need to be logged in and own the file to delete it
func handleDeleteFile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Delete file request received")
	
	// get the file path from the URL
	path := strings.TrimPrefix(r.URL.Path, "/files/")
	if path == "" {
		log.Println("DEBUG: Delete failed - no file ID provided")
		http.Error(w, "file ID required", http.StatusBadRequest)
		return
	}

	// make sure the user is logged in
	email, err := svc.Authorize(r)
	if err != nil {
		log.Printf("DEBUG: Delete unauthorized: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get the user info
	user, err := svc.GetUserByEmail(email)
	if err != nil {
		log.Printf("DEBUG: User lookup failed for delete: %v", err)
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Attempting to delete file %s for user %d", path, user.ID)

	// actually delete the file from my storage
	if err := svc.DeleteFile(context.Background(), user.ID, path); err != nil {
		log.Printf("DEBUG: Delete failed for user %d, file %s: %v", user.ID, path, err)
		http.Error(w, "could not delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Successfully deleted file %s for user %d", path, user.ID)
	w.WriteHeader(http.StatusNoContent)
}

// handleGetFile gets details about a specific book file
// this might be used to show file info before downloading or reading
func handleGetFile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Get file details request received")
	
	// get the file path from the URL
	path := strings.TrimPrefix(r.URL.Path, "/files/")
	if path == "" {
		log.Println("DEBUG: Get file failed - no file ID provided")
		http.Error(w, "file ID required", http.StatusBadRequest)
		return
	}

	// make sure the user is logged in
	email, err := svc.Authorize(r)
	if err != nil {
		log.Printf("DEBUG: Get file unauthorized: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get the user info
	user, err := svc.GetUserByEmail(email)
	if err != nil {
		log.Printf("DEBUG: User lookup failed for get file: %v", err)
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Getting file info for %s, user %d", path, user.ID)

	// get the file details from my storage
	fileInfo, err := svc.GetFileInfo(context.Background(), user.ID, path)
	if err != nil {
		log.Printf("DEBUG: File not found: %s for user %d", path, user.ID)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	log.Printf("DEBUG: Sending file info for %s", path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfo)
}

// ListFilesHandler is basically the same as the library handler
// I might remove this later since it's redundant
func ListFilesHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DEBUG: List files request received")
		
		if r.Method != http.MethodGet {
			log.Printf("DEBUG: List files method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handleGetLibrary(svc, w, r)
	}
}

// ProfileHandler handles user profile operations
// right now it supports GET (view profile) and PUT (update profile)
func ProfileHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Profile request: %s %s", r.Method, r.URL.Path)
		
		switch r.Method {
		case http.MethodGet:
			handleGetProfile(svc, w, r)
		case http.MethodPut:
			handleUpdateProfile(svc, w, r)
		default:
			log.Printf("DEBUG: Profile method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleGetProfile returns the user's profile information
// I only send back safe data, not passwords or sensitive stuff
func handleGetProfile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Getting user profile")
	
	// get the user from the JWT token
	email, err := svc.Authorize(r)
	if err != nil {
		log.Printf("DEBUG: Profile access unauthorized: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// look up the user in my database
	user, err := svc.GetUserByEmail(email)
	if err != nil {
		log.Printf("DEBUG: User lookup failed for profile: %v", err)
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Sending profile for user %d", user.ID)

	// only send back safe profile data, not sensitive stuff
	profile := map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"provider": user.Provider,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// handleUpdateProfile would update user profile information
// I haven't implemented this yet, but I'll probably need it later
func handleUpdateProfile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Profile update requested (not implemented yet)")
	// TODO: Implement profile updates
	w.WriteHeader(http.StatusNotImplemented)
}

// FileProxyHandler serves the actual book files to users
// this is important because it checks permissions before serving files
// users can only access their own books, not other people's
func FileProxyHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: File proxy request: %s", r.URL.Path)
		
		// only allow GET requests for file downloads
		if r.Method != http.MethodGet {
			log.Printf("DEBUG: File proxy method not allowed: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// make sure the user is logged in
		email, err := svc.Authorize(r)
		if err != nil {
			log.Printf("DEBUG: File proxy unauthorized: %v", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// get the user info
		user, err := svc.GetUserByEmail(email)
		if err != nil {
			log.Printf("DEBUG: User lookup failed for file proxy: %v", err)
			http.Error(w, "user lookup failed", http.StatusInternalServerError)
			return
		}

		// parse the file path from the URL
		// format is: /files/user-123/filename.epub
		path := strings.TrimPrefix(r.URL.Path, "/files/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			log.Printf("DEBUG: Invalid file path: %s", path)
			http.Error(w, "invalid file path", http.StatusBadRequest)
			return
		}

		bucket := parts[0]
		objectName := parts[1]

		// make sure the user is trying to access their own files
		// this is a security check to prevent accessing other people's books
		expectedBucket := fmt.Sprintf("user-%d", user.ID)
		if bucket != expectedBucket {
			log.Printf("DEBUG: Access denied - user %d tried to access bucket %s", user.ID, bucket)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		log.Printf("DEBUG: Serving file %s for user %d", objectName, user.ID)

		// get the file stream from my storage (MinIO)
		obj, err := svc.GetFileStream(r.Context(), user.ID, objectName)
		if err != nil {
			log.Printf("DEBUG: File not found: %s for user %d", objectName, user.ID)
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		defer obj.Close()

		// set headers for the file download
		w.Header().Set("Content-Type", "application/epub+zip")
		w.Header().Set("Cache-Control", "private, max-age=3600")

		// stream the file content to the user
		log.Printf("DEBUG: Streaming file %s to user %d", objectName, user.ID)
		io.Copy(w, obj)
	}
}
