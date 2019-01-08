package jwt_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/glynternet/mon/internal/jwt"
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

type mockHandler struct {
	request *http.Request
	writer  http.ResponseWriter
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.writer = w
	h.request = r
}

func TestAuthoriseClaims(t *testing.T) {
	t.Run("extracting error", func(t *testing.T) {
		authoriser := mockClaimsAuthorise{newClaims: &struct{}{}}
		extractor := mockClaimsExtractor{
			err: errors.New("extracting error"),
		}

		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		var next mockHandler
		handlerFn := jwt.AuthoriseClaims(logger, &extractor, &authoriser, &next)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)
		token := &josejwt.JSONWebToken{}
		handlerFn(token, w, r)

		assert.Equal(t, r, extractor.request)
		assert.Equal(t, token, extractor.token)
		assert.Equal(t, []interface{}{authoriser.newClaims}, extractor.values)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Unable to extract claims: extracting error\n", buf.String())
		assert.Nil(t, next.writer)
		assert.Nil(t, next.request)
	})

	t.Run("authorise error", func(t *testing.T) {
		authoriser := mockClaimsAuthorise{
			newClaims: &struct{}{},
			err:       errors.New("authorising error"),
		}

		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		var next mockHandler
		handlerFn := jwt.AuthoriseClaims(logger, &mockClaimsExtractor{}, &authoriser, &next)

		w := httptest.NewRecorder()
		handlerFn(nil, w, nil)

		assert.Equal(t, authoriser.newClaims, authoriser.authorisedClaims)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "Unauthorised token: authorising error\n", buf.String())
		assert.Nil(t, next.writer)
		assert.Nil(t, next.request)
	})

	t.Run("okay", func(t *testing.T) {
		next := mockHandler{}
		handlerFn := jwt.AuthoriseClaims(nil, &mockClaimsExtractor{}, &mockClaimsAuthorise{}, &next)

		var w http.ResponseWriter
		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)
		handlerFn(nil, w, r)
		assert.Equal(t, next.request, r)
		assert.Equal(t, next.writer, w)
	})
}
