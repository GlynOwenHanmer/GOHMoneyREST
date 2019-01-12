package monauth

import (
	"context"

	"golang.org/x/oauth2"
)

// AuthCodeExchanger provides a URL that the user can browse to to login and
// exchanges the recieved code for an oauth2.Token for authentication.
type AuthCodeExchanger interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}
