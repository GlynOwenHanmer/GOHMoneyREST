package monauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/mon/internal/router"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

// ensure that loginURLHandler is an AppJSONHandler
var _ router.AppJSONHandler = (&env{}).loginURLHandler
var _ router.AppJSONHandler = (&env{}).loginCallbackHandler

func Test_loginURLHandler(t *testing.T) {
	t.Run("generate state error", func(t *testing.T) {
		env := &env{
			generateStateValue: func() (string, error) {
				return "", errors.New("gen error")
			},
		}
		code, bod, err := env.loginURLHandler(nil)
		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), bod)
		assert.Error(t, err)
	})

	t.Run("response generated", func(t *testing.T) {
		exchanger := &mockAuthCodeExchanger{
			url: "hiyer",
		}
		e := env{
			generateStateValue: func() (string, error) {
				return "yo", nil
			},
			state:             state{},
			AuthCodeExchanger: exchanger,
		}
		code, bod, err := e.loginURLHandler(nil)
		assert.Equal(t, "yo", e.state.Get())
		assert.Nil(t, exchanger.options)
		assert.Equal(t, LoginURLResponse{LoginURL: "hiyer"}, bod)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, code)
	})
}

func Test_loginCallbackHandler(t *testing.T) {
	t.Run("expected state not set", func(t *testing.T) {
		env := &env{}
		code, bod, err := env.loginCallbackHandler(nil)
		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), bod)
		assert.Error(t, err)
	})

	t.Run("unexpected state", func(t *testing.T) {
		env := &env{
			state: state{state: "hiyer"},
		}
		r := httptest.NewRequest("ANY", "any://any", nil)
		code, bod, err := env.loginCallbackHandler(r)
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, "invalid state parameter", bod)
		assert.Error(t, err)
	})

	t.Run("no code parameter", func(t *testing.T) {
		state := state{state: "hiyer"}
		env := &env{state: state}
		r := httptest.NewRequest("ANY", "any://any", nil)
		vs := r.URL.Query()
		vs.Set("state", "hiyer")
		r.URL.RawQuery = vs.Encode()

		code, bod, err := env.loginCallbackHandler(r)
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, "code parameter not set", bod)
		assert.NoError(t, err)
	})

	t.Run("code exchange error", func(t *testing.T) {
		state := state{state: "hiyer"}
		exchanger := &mockAuthCodeExchanger{
			error: errors.New("code exchange error"),
		}
		env := &env{
			state:             state,
			AuthCodeExchanger: exchanger,
		}

		r := httptest.NewRequest("ANY", "any://any", nil)
		vs := r.URL.Query()
		vs.Set("state", "hiyer")
		vs.Set("code", "yo")
		r.URL.RawQuery = vs.Encode()

		code, bod, err := env.loginCallbackHandler(r)
		assert.Equal(t, exchanger.code, "yo")
		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), bod)
		assert.Equal(t, exchanger.error, errors.Cause(err))
	})

	t.Run("no id_token", func(t *testing.T) {
		exchanger := &mockAuthCodeExchanger{
			token: &oauth2.Token{},
		}
		env := &env{
			state:             state{state: "hiyer"},
			AuthCodeExchanger: exchanger,
		}

		r := httptest.NewRequest("ANY", "any://any", nil)
		vs := r.URL.Query()
		vs.Set("state", "hiyer")
		vs.Set("code", "yo")
		r.URL.RawQuery = vs.Encode()
		code, bod, err := env.loginCallbackHandler(r)
		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), bod)
		assert.Error(t, err)
	})

	t.Run("id_token is empty", func(t *testing.T) {
		exchanger := &mockAuthCodeExchanger{
			token: (&oauth2.Token{}).WithExtra(
				map[string]interface{}{"id_token": ""},
			),
		}
		env := &env{
			state:             state{state: "hiyer"},
			AuthCodeExchanger: exchanger,
		}

		r := httptest.NewRequest("ANY", "any://any", nil)
		vs := r.URL.Query()
		vs.Set("state", "hiyer")
		vs.Set("code", "yo")
		r.URL.RawQuery = vs.Encode()
		code, bod, err := env.loginCallbackHandler(r)
		assert.Equal(t, http.StatusInternalServerError, code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), bod)
		assert.Error(t, err)
	})

	t.Run("all good", func(t *testing.T) {
		exchanger := &mockAuthCodeExchanger{
			token: (&oauth2.Token{}).WithExtra(
				map[string]interface{}{"id_token": "woooooh"},
			),
		}
		env := &env{
			state:             state{state: "hiyer"},
			AuthCodeExchanger: exchanger,
		}

		r := httptest.NewRequest("ANY", "any://any", nil)
		vs := r.URL.Query()
		vs.Set("state", "hiyer")
		vs.Set("code", "yo")
		r.URL.RawQuery = vs.Encode()
		code, bod, err := env.loginCallbackHandler(r)
		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, http.StatusText(http.StatusOK), bod)
		assert.NoError(t, err)
	})
}

type mockAuthCodeExchanger struct {
	url     string
	state   string
	options []oauth2.AuthCodeOption

	code  string
	token *oauth2.Token
	error
}

func (e *mockAuthCodeExchanger) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	e.state = state
	e.options = opts
	return e.url
}

func (e *mockAuthCodeExchanger) Exchange(_ context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	e.code = code
	return e.token, e.error
}
