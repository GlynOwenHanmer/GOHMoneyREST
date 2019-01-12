package monauth

import "net/http"

// env holds all items required to process the login flow
type env struct {
	AuthCodeExchanger
	generateStateValue func() (string, error)
	state
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
	e.Set(state)
	return http.StatusOK, LoginURLResponse{LoginURL: url}, nil
}

// LoginURLResponse is the response body that is returned when requesting a new
// login URL.
type LoginURLResponse struct {
	LoginURL string `json:"login_url"`
}
