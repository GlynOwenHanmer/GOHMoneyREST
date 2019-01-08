package jwt

import (
	"log"
	"net/http"

	"gopkg.in/square/go-jose.v2/jwt"
)

// Validate creates an http.HandlerFunc that will validate a request and its
// JWT.  On successful validation, the token will be passed to the given
// TokenHandlerFunc.
// The generated HandlerFunc will panic if either the given RequestValidator or
// next http.Handler are nil.
func Validate(logger *log.Logger, rv RequestValidator, next TokenHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := rv.ValidateRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			if logger != nil {
				logger.Printf("Unauthorised request: %+v", err)
			}
			return
		}
		next(t, w, r)
	}
}

// RequestValidator is used to validate a given request and its JWT
type RequestValidator interface {
	// ValidateRequest validates a given request, returning the validated JSONWebToken.
	ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error)
}

// TokenHandlerFunc is handler that will use a token within its logic
type TokenHandlerFunc func(*jwt.JSONWebToken, http.ResponseWriter, *http.Request)
