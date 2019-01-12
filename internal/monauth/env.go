package monauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// env holds all items required to process the login flow
type env struct {
	AuthCodeExchanger
	generateStateValue func() (string, error)
	state
	tokenExchangeTimout time.Duration
}

// loginURLHandler generates a login URL that the user can navigate to to login
func (e *env) loginURLHandler(_ *http.Request) (int, interface{}, error) {
	state, err := e.generateStateValue()
	if err != nil {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err
	}
	url := e.AuthCodeURL(state)
	e.state.Set(state)
	return http.StatusOK, LoginURLResponse{LoginURL: url}, nil
}

// LoginURLResponse is the response body that is returned when requesting a new
// login URL.
type LoginURLResponse struct {
	LoginURL string `json:"login_url"`
}

// loginCallbackHandler exchanges the query paramater "code" for a jwt id token.
func (e *env) loginCallbackHandler(r *http.Request) (int, interface{}, error) {
	expectedState := e.state.Get()
	if strings.TrimSpace(expectedState) == "" {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			errors.New("expected state is not set")
	}

	state := r.URL.Query().Get("state")
	if expectedState != state {
		return http.StatusBadRequest,
			"invalid state parameter",
			fmt.Errorf("expected state %q but got %q", expectedState, state)
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		return http.StatusBadRequest, "code parameter not set", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.tokenExchangeTimout)
	defer cancel()
	token, err := e.Exchange(ctx, code)
	if err != nil {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			errors.Wrap(err, "exchanging code for token")
	}

	jwt, ok := token.Extra("id_token").(string)
	if !ok {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			errors.New("id_token not present")
	}

	if jwt == "" {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			errors.New("id_token is empty")
	}

	return http.StatusOK, http.StatusText(http.StatusOK), nil
}
