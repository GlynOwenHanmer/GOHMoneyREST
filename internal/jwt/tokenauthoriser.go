package jwt

import (
	"net/http"

	"gopkg.in/square/go-jose.v2/jwt"
)

// TokenAuthoriser authorises a token and request
type TokenAuthoriser interface {
	Authorise(*jwt.JSONWebToken, *http.Request) error
}

// TokenAuthoriserFunc is a function authorises a token and request
type TokenAuthoriserFunc func(*jwt.JSONWebToken, *http.Request) error

// Authorise ensures that TokenAuthoriserFunc satisfies the TokenAuthoriser
// interface
func (fn TokenAuthoriserFunc) Authorise(t *jwt.JSONWebToken, r *http.Request) error {
	return fn(t, r)
}
