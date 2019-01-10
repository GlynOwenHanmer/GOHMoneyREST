package jwt

import (
	"net/http"

	"github.com/pkg/errors"
)

// RequestTokenAuthoriser validates
type RequestTokenAuthoriser struct {
	RequestValidator
	TokenAuthoriser
}

// Authorise validates a request's JWT and checks that the JWT is authorised,
// returning an error for an unauthorised token.
func (ra RequestTokenAuthoriser) Authorise(r *http.Request) error {
	token, err := ra.ValidateRequest(r)
	if err != nil {
		return errors.Wrap(err, "validating request")
	}
	return ra.TokenAuthoriser.Authorise(token, r)
}
