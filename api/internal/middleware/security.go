package middleware

import (
	"net/http"
	"strings"

	"github.com/KaikSelhorst/shortener/internal/httputil"
)

// SecurityHeaders sets defensive HTTP headers on every response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		next.ServeHTTP(w, r)
	})
}

// RequireJSON rejects POST/PUT/PATCH requests that are not application/json.
func RequireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				httputil.WriteError(w, http.StatusUnsupportedMediaType, "content-type must be application/json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
