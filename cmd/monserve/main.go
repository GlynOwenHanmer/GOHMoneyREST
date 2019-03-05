package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/glynternet/mon/internal/auth/auth0"
	"github.com/glynternet/mon/internal/router"
	"github.com/glynternet/mon/internal/versioncmd"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "monserve"

	// viper keys
	keyPort            = "port"
	keySSLCertificate  = "ssl-certificate"
	keySSLKey          = "ssl-key"
	keyAuth0Domain     = "auth0-domain"
	keyAuthorisedEmail = "authorised-email"
	keyDBHost          = "db-host"
	keyDBUser          = "db-user"
	keyDBPassword      = "db-password"
	keyDBName          = "db-name"
	keyDBSSLMode       = "db-sslmode"
)

// to be changed using ldflags with the go build command
var version = "unknown"

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	var cmdServe = &cobra.Command{
		Use: appName,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := newStorage(
				viper.GetString(keyDBHost),
				viper.GetString(keyDBUser),
				viper.GetString(keyDBPassword),
				viper.GetString(keyDBName),
				viper.GetString(keyDBSSLMode),
			)
			if err != nil {
				return errors.Wrap(err, "error creating storage")
			}
			var handler http.Handler
			handler, err = router.New(store, logger)
			if err != nil {
				return errors.Wrap(err, "error creating new server")
			}

			auth0Domain := viper.GetString(keyAuth0Domain)
			if auth0Domain == "" {
				logger.Print("NO AUTH0 DOMAIN GIVEN.\nNO AUTHORISATION WILL BE USED!")
			} else {
				handler, err = withAuth(
					logger,
					auth0Domain,
					viper.GetString(keyAuthorisedEmail),
					handler,
				)
				if err != nil {
					return errors.Wrap(err, "setting up auth")
				}
			}

			serveFn := newServeFn(
				logger,
				viper.GetString(keySSLCertificate),
				viper.GetString(keySSLKey),
			)
			addr := ":" + viper.GetString(keyPort)
			logger.Printf("Serving at %s", addr)
			return serveFn(addr, handler)
		},
	}

	cmdServe.AddCommand(versioncmd.New(version, os.Stdout))

	cobra.OnInitialize(viperAutoEnvVar)
	cmdServe.Flags().String(keyPort, "80", "server listening port")
	cmdServe.Flags().String(keySSLCertificate, "", "path to SSL certificate, leave empty for http")
	cmdServe.Flags().String(keySSLKey, "", "path to SSL key, leave empty for https")
	cmdServe.Flags().String(keyAuth0Domain, "",
		"auth0 domain to use for authentication. If left empty, no authorisation will be set",
	)
	cmdServe.Flags().String(keyAuthorisedEmail, "",
		fmt.Sprintf(
			"email address of authorised user. Must be set if %s is set. Will have no effect if %s is not set",
			keyAuth0Domain, keyAuth0Domain),
	)
	cmdServe.Flags().String(keyDBHost, "", "host address of the DB backend")
	cmdServe.Flags().String(keyDBName, "", "name of the DB set to use")
	cmdServe.Flags().String(keyDBUser, "", "DB user to authenticate with")
	cmdServe.Flags().String(keyDBPassword, "", "DB password to authenticate with")
	cmdServe.Flags().String(keyDBSSLMode, "", "DB SSL mode to use")
	err := viper.BindPFlags(cmdServe.Flags())
	if err != nil {
		logger.Printf("unable to BindPFlags: %v", err)
		os.Exit(1)
	}

	if err := cmdServe.Execute(); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}

// withAuth wraps the given handler in an authorisation handler, returning the wrapped handler
func withAuth(logger *log.Logger, auth0Domain, email string, handler http.Handler) (http.Handler, error) {
	authoriser, err := auth0.UserEmailAuthoriser(auth0Domain, email)
	if err != nil {
		return nil, errors.Wrap(err, "creating authorisation middleware")
	}
	logger.Printf("Auth0 domain: %s", auth0Domain)
	return router.WithAuthoriser(logger, authoriser, handler), nil
}

func viperAutoEnvVar() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func newStorage(host, user, password, dbname, sslmode string) (storage.Storage, error) {
	cs, err := postgres.NewConnectionString(host, user, password, dbname, sslmode)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection string: %v", err)
	}
	return postgres.New(cs)
}

// newServeFn returns a function that can be used to start a server.
// newServeFn will provide an HTTPS server if either the given certPath or
// keyPath are non-empty, otherwise newServeFn will provide an HTTP server.
func newServeFn(logger *log.Logger, certPath, keyPath string) func(string, http.Handler) error {
	if len(certPath) == 0 && len(keyPath) == 0 {
		logger.Printf("Using HTTP")
		return http.ListenAndServe
	}
	logger.Printf("Using HTTPS")
	return func(addr string, handler http.Handler) error {
		return http.ListenAndServeTLS(addr, certPath, keyPath, handler)
	}
}
