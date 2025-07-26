// internal/handlers/library.go
package handlers

import (
	"bookapp/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type LibraryResponse struct {
	Files []FileInfo `json:"files"`
	Total int        `json:"total"`
}

type FileInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Size     int64  `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

func LibraryHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetLibrary(svc, w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// internal/handlers/library.go - Updated handleGetLibrary function
func handleGetLibrary(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	// Extract user from context (set by JWT middleware)
	email, err := svc.Authorize(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := svc.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	urls, err := svc.ListFiles(context.Background(), user.ID)
	if err != nil {
		http.Error(w, "could not list files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert URLs to FileInfo structs with proxy URLs
	files := make([]FileInfo, len(urls))
	for i, url := range urls {
		// Extract filename from MinIO URL
		parts := strings.Split(url, "/")
		filename := "unknown"
		if len(parts) > 0 {
			filename = parts[len(parts)-1]
		}

		// Create proxy URL instead of direct MinIO URL
		proxyURL := fmt.Sprintf("/api/protected/files/user-%d/%s", user.ID, filename)

		files[i] = FileInfo{
			ID:   strconv.Itoa(i), // Consider using actual file IDs
			Name: filename,
			URL:  proxyURL, // Use proxy URL instead of direct MinIO URL
		}
	}

	response := LibraryResponse{
		Files: files,
		Total: len(files),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Additional handler for file operations
func FileHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handleDeleteFile(svc, w, r)
		case http.MethodGet:
			handleGetFile(svc, w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleDeleteFile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	// Extract file ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/files/")
	if path == "" {
		http.Error(w, "file ID required", http.StatusBadRequest)
		return
	}

	// Get user
	email, err := svc.Authorize(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := svc.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	// TODO: Implement delete file in service
	if err := svc.DeleteFile(context.Background(), user.ID, path); err != nil {
		http.Error(w, "could not delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGetFile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	// Extract file ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/files/")
	if path == "" {
		http.Error(w, "file ID required", http.StatusBadRequest)
		return
	}

	// Get user
	email, err := svc.Authorize(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := svc.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	// TODO: Implement get file details in service
	fileInfo, err := svc.GetFileInfo(context.Background(), user.ID, path)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfo)
}

// List files handler (separate from library for more specific file operations)
func ListFilesHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handleGetLibrary(svc, w, r)
	}
}

// Profile handler for user data
func ProfileHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetProfile(svc, w, r)
		case http.MethodPut:
			handleUpdateProfile(svc, w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetProfile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	email, err := svc.Authorize(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := svc.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "user lookup failed", http.StatusInternalServerError)
		return
	}

	// Don't expose sensitive data
	profile := map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"provider": user.Provider,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func handleUpdateProfile(svc *auth.Service, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement profile updates
	w.WriteHeader(http.StatusNotImplemented)
}

// Add this to handlers/library.go

func FileProxyHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Authorize user
		email, err := svc.Authorize(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := svc.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "user lookup failed", http.StatusInternalServerError)
			return
		}

		// Extract bucket and object from path: /files/user-123/filename.epub
		path := strings.TrimPrefix(r.URL.Path, "/files/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			http.Error(w, "invalid file path", http.StatusBadRequest)
			return
		}

		bucket := parts[0]
		objectName := parts[1]

		// Verify user owns this bucket
		expectedBucket := fmt.Sprintf("user-%d", user.ID)
		if bucket != expectedBucket {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		// Get file from MinIO and stream it
		obj, err := svc.GetFileStream(r.Context(), user.ID, objectName)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		defer obj.Close()

		// Set appropriate headers
		w.Header().Set("Content-Type", "application/epub+zip")
		w.Header().Set("Cache-Control", "private, max-age=3600")

		// Stream the file
		io.Copy(w, obj)
	}
}
