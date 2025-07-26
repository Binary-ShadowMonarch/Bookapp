package middleware

import (
	"log"
	"net/http"
)

// allowedOrigins is your “allow list” of dev domains.
var allowedOrigins = map[string]bool{
	"http://localhost:5173":             true,
	"http://localhost:3000":             true,
	"http://localhost:4353":             true,
	"https://books.saurabpoudel.com.np": true,
}

// CORS wraps your handler and enforces origin checks.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Print(origin)
		if origin == "" {
			// not a CORS request—proceed
			next.ServeHTTP(w, r)
			return
		}

		// if not in our allow‑list, reject immediately
		if !allowedOrigins[origin] {
			http.Error(w, "forbidden origin", http.StatusForbidden)
			return
		}

		// safe to echo back the origin
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin") // tells caches this header varies by Origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
