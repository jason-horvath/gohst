package middleware

import (
	"net/http"
)

func Template(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// sess, _ := session.SM.GetSession(r)


		// Add template logic here
		// For now, just pass the request
		next.ServeHTTP(w, r)
	})
}
