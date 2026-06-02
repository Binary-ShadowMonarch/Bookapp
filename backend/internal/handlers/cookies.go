package handlers

import (
	"net/http"
	"os"
	"strings"
)

// secureCookie returns whether auth cookies should be marked Secure.
// It prefers an explicit env override, then falls back to request context.
func secureCookie(r *http.Request) bool {
	if v := strings.TrimSpace(os.Getenv("COOKIE_SECURE")); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}

	if r.TLS != nil {
		return true
	}
	if strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		return true
	}
	return false
}
