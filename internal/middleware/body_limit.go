package middleware

import (
	"net/http"
)

// MaxBodySize limits the request body size. Requests exceeding the limit
// receive a 413 Request Entity Too Large response.
// Use 64*1024 (64 KB) for auth forms, 1<<20 (1 MB) for general forms.
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}
