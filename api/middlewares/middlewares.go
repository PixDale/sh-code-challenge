package middlewares

import (
	"errors"
	"net/http"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/responses"
)

// SetMiddlewareJSON defines a middleware to set the content type of a request to JSON
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication defines a middleware to verify the integrity of the token within the request
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
