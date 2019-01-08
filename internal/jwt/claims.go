package jwt

import (
	"log"
	"net/http"

	"gopkg.in/square/go-jose.v2/jwt"
)

// AuthoriseClaims extracts the claims from a given JWT and authorises a request against the claims.
// If the claims are authorised, the next http.Handler will be called, otherwise an Unauthorised status code will be written to the response.
// It is the responsibility of the caller to provide non-nil values
func AuthoriseClaims(logger *log.Logger, ce ClaimsExtractor, ca ClaimsAuthoriser, next http.Handler) TokenHandlerFunc {
	return func(token *jwt.JSONWebToken, w http.ResponseWriter, r *http.Request) {
		claims := ca.NewClaims()
		err := ce.Claims(r, token, claims)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			if logger != nil {
				logger.Printf("Unable to extract claims: %+v", err)
			}
			return
		}

		err = ca.Authorise(claims)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			if logger != nil {
				logger.Printf("Unauthorised token: %+v", err)
			}
			return
		}

		next.ServeHTTP(w, r)
	}
}

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
