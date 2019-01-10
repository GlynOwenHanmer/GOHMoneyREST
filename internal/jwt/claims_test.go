package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/mon/internal/jwt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	josejwt "gopkg.in/square/go-jose.v2/jwt"
)

type mockClaimsExtractor struct {
	err     error
	request *http.Request
	token   *josejwt.JSONWebToken
	values  []interface{}
}

func (e *mockClaimsExtractor) Claims(r *http.Request, token *josejwt.JSONWebToken, values ...interface{}) error {
	e.request = r
	e.token = token
	e.values = values
	return e.err
}

type mockClaimsAuthorise struct {
	err              error
	newClaims        interface{}
	authorisedClaims interface{}
}

func (p mockClaimsAuthorise) NewClaims() interface{} {
	return p.newClaims
}

func (p *mockClaimsAuthorise) Authorise(claims interface{}) error {
	p.authorisedClaims = claims
	return p.err
}

func TestAuthoriseClaims(t *testing.T) {
	t.Run("extracting error", func(t *testing.T) {
		authoriser := mockClaimsAuthorise{newClaims: &struct{}{}}
		extractor := mockClaimsExtractor{
			err: errors.New("extracting error"),
		}

		handlerFn := jwt.AuthoriseClaims(&extractor, &authoriser)

		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)
		token := &josejwt.JSONWebToken{}
		err := handlerFn(token, r)

		assert.Error(t, err)
		assert.Equal(t, extractor.err, errors.Cause(err))
		assert.Equal(t, r, extractor.request)
		assert.Equal(t, token, extractor.token)
		assert.Equal(t, []interface{}{authoriser.newClaims}, extractor.values)
	})

	t.Run("authorise error", func(t *testing.T) {
		authoriser := mockClaimsAuthorise{
			newClaims: &struct{}{},
			err:       errors.New("authorising error"),
		}

		handlerFn := jwt.AuthoriseClaims(&mockClaimsExtractor{}, &authoriser)
		err := handlerFn(nil, nil)

		assert.Error(t, err)
		assert.Equal(t, authoriser.newClaims, authoriser.authorisedClaims)
	})

	t.Run("okay", func(t *testing.T) {
		handlerFn := jwt.AuthoriseClaims(&mockClaimsExtractor{}, &mockClaimsAuthorise{})

		err := handlerFn(nil, nil)
		assert.Nil(t, err)
	})
}
