package jwt_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/mon/internal/jwt"
	"github.com/stretchr/testify/assert"
	josejwt "gopkg.in/square/go-jose.v2/jwt"
)

type mockTokenHandler struct {
	called  bool
	request *http.Request
	token   *josejwt.JSONWebToken
}

func (h *mockTokenHandler) ServeHTTP(t *josejwt.JSONWebToken, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	h.called = true
	h.request = r
	h.token = t
}

type mockRequestValidator struct {
	token   *josejwt.JSONWebToken
	err     error
	request *http.Request
}

func (rv *mockRequestValidator) ValidateRequest(r *http.Request) (*josejwt.JSONWebToken, error) {
	rv.request = r
	return rv.token, rv.err
}

func TestValidate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		requestValidator := mockRequestValidator{
			err: errors.New("invalid request"),
		}
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		var next mockTokenHandler
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)

		hf := jwt.Validate(logger, &requestValidator, (&next).ServeHTTP)
		hf(w, r)

		assert.Equal(t, "Unauthorised request: invalid request\n", buf.String())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, next.called)
	})

	t.Run("valid", func(t *testing.T) {
		requestValidator := mockRequestValidator{
			token: &josejwt.JSONWebToken{},
		}
		var buf bytes.Buffer
		logger := log.New(&buf, "", 0)
		var next mockTokenHandler
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://any/", nil)

		hf := jwt.Validate(logger, &requestValidator, (&next).ServeHTTP)
		hf(w, r)

		assert.Equal(t, "", buf.String())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, true, next.called)
		assert.Equal(t, r, next.request)
	})
}
