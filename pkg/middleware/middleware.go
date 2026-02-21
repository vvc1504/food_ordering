package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// APIKeyAuth validates the api_key header.
func APIKeyAuth(apiKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("api_key")
		if val != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Logger logs the details of each request.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Dur("duration", time.Since(start)).
			Msg("Request handled")
	}
}
