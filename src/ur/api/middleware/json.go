package middleware

import (
	"net/http"
)

func ContentTypeJSON(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		inner.ServeHTTP(w, r)
	})
}
