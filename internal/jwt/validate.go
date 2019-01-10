package jwt

import (
	"net/http"

	"gopkg.in/square/go-jose.v2/jwt"
)

// RequestValidator is used to validate a given request and its JWT
type RequestValidator interface {
	// ValidateRequest validates a given request, returning the validated JSONWebToken.
	ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error)
}
