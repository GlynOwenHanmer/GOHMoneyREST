package jwt

import (
	"net/http"

	"gopkg.in/square/go-jose.v2/jwt"
)

// TokenAuthoriser authorises a token and request
type TokenAuthoriser interface {
	Authorise(*jwt.JSONWebToken, *http.Request) error
}
