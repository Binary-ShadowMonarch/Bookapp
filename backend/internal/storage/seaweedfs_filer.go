package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type SeaweedFSConfig struct {
	FilerURL string
	BasePath string
	Timeout  time.Duration
}

type SeaweedFSFiler struct {
	baseURL  string
	basePath string
	client   *http.Client
}

func NewSeaweedFSFiler(cfg SeaweedFSConfig) (*SeaweedFSFiler, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.FilerURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("seaweedfs filer url is required")
	}
	basePath := normalizeBasePath(cfg.BasePath)
	client := &http.Client{Timeout: cfg.Timeout}
	if cfg.Timeout == 0 {
		client.Timeout = 30 * time.Second
	}

	return &SeaweedFSFiler{
		baseURL:  baseURL,
		basePath: basePath,
		client:   client,
	}, nil
}

func normalizeBasePath(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return "/bookapp"
	}
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	return strings.TrimRight(base, "/")
}

func (s *SeaweedFSFiler) EnsureScope(ctx context.Context, scope Scope) error {
	// Direct filer operations create directories on demand, so this is a no-op for now.
	return ValidateScope(scope)
}

func (s *SeaweedFSFiler) Put(ctx context.Context, scope Scope, objectName string, reader io.Reader, size int64, contentType string) error {
	objectPath, err := s.objectPath(scope, objectName)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.urlForPath(objectPath), reader)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if size >= 0 {
		req.ContentLength = size
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("seaweedfs put failed: %s", readErrorBody(resp.Body, resp.Status))
	}
	return nil
}

func (s *SeaweedFSFiler) Get(ctx context.Context, scope Scope, objectName string) (*ObjectReader, error) {
	objectPath, err := s.objectPath(scope, objectName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.urlForPath(objectPath), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("seaweedfs get failed: %s", readErrorBody(resp.Body, resp.Status))
	}

	return &ObjectReader{
		Body:        resp.Body,
		ContentType: resp.Header.Get("Content-Type"),
		Size:        resp.ContentLength,
	}, nil
}

func (s *SeaweedFSFiler) List(ctx context.Context, scope Scope) ([]ObjectInfo, error) {
	prefix, err := s.scopePrefix(scope)
	if err != nil {
		return nil, err
	}

	directoryURL := s.urlForPath(strings.TrimSuffix(prefix, "/") + "/")
	url := directoryURL + "?limit=10000"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("seaweedfs list failed: %s", readErrorBody(resp.Body, resp.Status))
	}

	var listing filerListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return nil, fmt.Errorf("decoding seaweedfs list response: %w", err)
	}

	objects := make([]ObjectInfo, 0, len(listing.Entries))
	for _, entry := range listing.Entries {
		if entry.IsDirectory {
			continue
		}
		name := entry.Name
		if name == "" && entry.FullPath != "" {
			name = path.Base(entry.FullPath)
		}
		if name == "" {
			continue
		}
		objects = append(objects, ObjectInfo{
			Key:          name,
			Size:         entry.Size,
			ContentType:  entry.Mime,
			LastModified: unixSecondsToTime(entry.Mtime),
		})
	}

	return objects, nil
}

func (s *SeaweedFSFiler) Stat(ctx context.Context, scope Scope, objectName string) (*ObjectInfo, error) {
	objectPath, err := s.objectPath(scope, objectName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, s.urlForPath(objectPath), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("seaweedfs stat failed: %s", readErrorBody(resp.Body, resp.Status))
	}

	info := &ObjectInfo{
		Key:         objectName,
		Size:        resp.ContentLength,
		ContentType: resp.Header.Get("Content-Type"),
	}
	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		if parsed, err := time.Parse(time.RFC1123, lastModified); err == nil {
			info.LastModified = parsed
		}
	}
	return info, nil
}

func (s *SeaweedFSFiler) Delete(ctx context.Context, scope Scope, objectName string) error {
	objectPath, err := s.objectPath(scope, objectName)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.urlForPath(objectPath), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("seaweedfs delete failed: %s", readErrorBody(resp.Body, resp.Status))
	}

	return nil
}

func (s *SeaweedFSFiler) scopePrefix(scope Scope) (string, error) {
	if err := ValidateScope(scope); err != nil {
		return "", err
	}
	if scope.Namespace == NamespaceUser {
		return path.Join(s.basePath, "users", strconv.Itoa(scope.UserID)), nil
	}
	return path.Join(s.basePath, "global"), nil
}

func (s *SeaweedFSFiler) objectPath(scope Scope, objectName string) (string, error) {
	if err := ValidateObjectName(objectName); err != nil {
		return "", err
	}
	prefix, err := s.scopePrefix(scope)
	if err != nil {
		return "", err
	}
	return path.Join(prefix, objectName), nil
}

func (s *SeaweedFSFiler) urlForPath(p string) string {
	if p == "" {
		return s.baseURL
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return s.baseURL + p
}

type filerListResponse struct {
	Entries []filerEntry `json:"Entries"`
}

type filerEntry struct {
	Name        string `json:"Name"`
	FullPath    string `json:"FullPath"`
	IsDirectory bool   `json:"IsDirectory"`
	Size        int64  `json:"Size"`
	Mime        string `json:"Mime"`
	Mtime       int64  `json:"Mtime"`
}

func unixSecondsToTime(value int64) time.Time {
	if value <= 0 {
		return time.Time{}
	}
	return time.Unix(value, 0).UTC()
}

func readErrorBody(body io.Reader, status string) string {
	const maxSize = 4096
	buf := make([]byte, maxSize)
	n, _ := body.Read(buf)
	payload := strings.TrimSpace(string(buf[:n]))
	if payload == "" {
		return status
	}
	return status + ": " + payload
}
