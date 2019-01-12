package monauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glynternet/mon/internal/monauth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestServeMux(t *testing.T) {
	exchanger := &mockAuthCodeExchanger{
		url: "hiyer",
		token: (&oauth2.Token{}).WithExtra(
			map[string]interface{}{"id_token": "woooooh"},
		),
	}
	m := monauth.ServeMux(exchanger, time.Second)
	assert.NotNil(t, m)

	r := httptest.NewRequest(http.MethodGet, "any://any/loginurl", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualURLResponse monauth.LoginURLResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualURLResponse)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	expectedURLResponse := monauth.LoginURLResponse{
		LoginURL: "hiyer",
	}
	assert.Equal(t, expectedURLResponse, actualURLResponse)
	r = httptest.NewRequest(http.MethodGet, "any://any/logincallback", nil)
	vs := r.URL.Query()
	vs.Set("state", exchanger.state)
	vs.Set("code", "any_code_anytime")
	r.URL.RawQuery = vs.Encode()

	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualCallbackResponse monauth.LoginCallbackResponse
	err = json.Unmarshal(w.Body.Bytes(), &actualCallbackResponse)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	expectedCallbackResponse := monauth.LoginCallbackResponse{
		Token: "woooooh",
	}
	assert.Equal(t, expectedCallbackResponse, actualCallbackResponse)
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
