package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2/jwt"
)

func TestRequestTokenAuthoriser_Authorise(t *testing.T) {
	t.Run("invalid request", func(t *testing.T) {
		validator := mockRequestValidator{
			err: errors.New("validation error"),
		}
		rta := RequestTokenAuthoriser{
			RequestValidator: &validator,
		}

		r := httptest.NewRequest(http.MethodGet, "http://boop", nil)
		err := rta.Authorise(r)
		assert.Equal(t, r, validator.request)
		assert.Equal(t, validator.err, errors.Cause(err))
	})

	t.Run("unauhorised token", func(t *testing.T) {
		validator := mockRequestValidator{
			token: &jwt.JSONWebToken{},
		}
		authoriser := mockTokenAuthoriser{
			err: errors.New("auth error"),
		}
		rta := RequestTokenAuthoriser{
			RequestValidator: &validator,
			TokenAuthoriser:  &authoriser,
		}

		r := httptest.NewRequest(http.MethodGet, "http://boop", nil)
		err := rta.Authorise(r)
		assert.Equal(t, r, authoriser.request)
		assert.Equal(t, validator.token, authoriser.token)
		assert.Equal(t, authoriser.err, err)
	})

	t.Run("okay", func(t *testing.T) {
		rta := RequestTokenAuthoriser{
			RequestValidator: &mockRequestValidator{},
			TokenAuthoriser:  &mockTokenAuthoriser{},
		}

		r := httptest.NewRequest(http.MethodGet, "http://boop", nil)
		err := rta.Authorise(r)
		assert.NoError(t, err)
	})
}

type mockRequestValidator struct {
	request *http.Request
	token   *jwt.JSONWebToken
	err     error
}

func (mrv *mockRequestValidator) ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error) {
	mrv.request = r
	return mrv.token, mrv.err
}

type mockTokenAuthoriser struct {
	token   *jwt.JSONWebToken
	request *http.Request
	err     error
}

func (mta *mockTokenAuthoriser) Authorise(t *jwt.JSONWebToken, r *http.Request) error {
	mta.request = r
	mta.token = t
	return mta.err
}
