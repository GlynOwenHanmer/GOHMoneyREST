package auth0

import (
	"errors"
	"fmt"

	"github.com/auth0-community/go-auth0"
	"github.com/glynternet/mon/internal/auth"
	"github.com/glynternet/mon/internal/jwt"
	"github.com/glynternet/mon/internal/model"
	"github.com/glynternet/mon/internal/router"
	"gopkg.in/square/go-jose.v2"
)

// UserEmailAuthoriser authorises requests against their contained JWT.
// Requests are deemed to be valid if they contain a JWT Bearer token with the
// given verified email address, signed by the given auth0domain.
func UserEmailAuthoriser(auth0Domain, email string) (router.RequestAuthoriser, error) {
	if auth0Domain == "" {
		return nil, errors.New("no auth0Domain given")
	}
	if email == "" {
		return nil, errors.New("no email given")
	}

	u := model.User{
		Email:         email,
		EmailVerified: true,
	}

	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain)
	auth0DomainURL := fmt.Sprintf("https://%s/", auth0Domain)

	jwkClient := auth0.NewJWKClient(
		auth0.JWKClientOptions{URI: jwksURL},
		auth0.RequestTokenExtractorFunc(auth0.FromHeader),
	)
	auth0Conf := auth0.NewConfiguration(jwkClient, []string{}, auth0DomainURL, jose.RS256)
	jwtValidator := auth0.NewValidator(auth0Conf, auth0.RequestTokenExtractorFunc(auth0.FromHeader))

	tokenAuth := jwt.AuthoriseClaims(
		jwtValidator,
		auth.UserClaimsAuthoriser(u),
	)
	return jwt.RequestTokenAuthoriser{
		RequestValidator: jwtValidator,
		TokenAuthoriser:  tokenAuth,
	}, nil
}
