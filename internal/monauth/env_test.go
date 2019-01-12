package monauth

import (
	"context"
	"net/http"
	"testing"

	"github.com/glynternet/mon/internal/router"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

// ensure that loginURLHandler is an AppJSONHandler
var _ router.AppJSONHandler = (&env{}).loginURLHandler

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

	t.Run("state becomes set", func(t *testing.T) {
		e := env{
			generateStateValue: func() (string, error) {
				return "yo", nil
			},
		}
		_, _, err := e.loginURLHandler(nil)
		assert.NoError(t, err)
		assert.Equal(t, "yo", e.state.Get())
	})

	t.Run("response generated", func(t *testing.T) {
		exchanger := &mockAuthCodeExchanger{
			url: "hiyer",
		}
		e := env{
			generateStateValue: func() (string, error) {
				return "yo", nil
			},
			AuthCodeExchanger: exchanger,
		}
		code, bod, err := e.loginURLHandler(nil)
		assert.Equal(t, "yo", exchanger.state)
		assert.Nil(t, exchanger.options)
		assert.Equal(t, LoginURLResponse{LoginURL: "hiyer"}, bod)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, code)
	})
}

type mockAuthCodeExchanger struct {
	url     string
	state   string
	options []oauth2.AuthCodeOption
}

func (e *mockAuthCodeExchanger) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	e.state = state
	e.options = opts
	return e.url
}

func (e *mockAuthCodeExchanger) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	panic("not implemented")
}
