package jwt

import (
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/jwt"
)

// AuthoriseClaims extracts the claims from a given JWT and authorises a
// request against the claims, returning an error for any unauthorised
// requests.
func AuthoriseClaims(ce ClaimsExtractor, ca ClaimsAuthoriser) TokenAuthoriserFunc {
	return func(token *jwt.JSONWebToken, r *http.Request) error {
		claims := ca.NewClaims()
		err := ce.Claims(r, token, claims)
		if err != nil {
			return errors.Wrap(err, "Unable to extract claims")
		}

		err = ca.Authorise(claims)
		if err != nil {
			return errors.Wrap(err, "Unauthorised token")
		}

		return nil
	}
}

// ClaimsAuthoriser generates a new claims object, which should have the token
// claims json unmarshalled into it, and authorises the claims, returning an
// error for unauthorised claims.
type ClaimsAuthoriser interface {
	// NewClaims provides an object that the claims of the JWT token can
	// be json unmarshalled into.
	NewClaims() interface{}
	// Authorise processes the claims, returning an error if they are unauthorised.
	Authorise(claims interface{}) error
}

// ClaimsExtractor extracts the claims from a JWT, possibly using other information from the given request.
// The claims are json unmarshalled into each of the given values.
type ClaimsExtractor interface {
	Claims(r *http.Request, token *jwt.JSONWebToken, values ...interface{}) error
}
