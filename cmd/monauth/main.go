package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/glynternet/mon/internal/middleware"
	"github.com/glynternet/mon/internal/monauth"
	"github.com/glynternet/mon/internal/versioncmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	appName = "monauth"

	// viper keys
	keyPort              = "port"
	keySSLCertificate    = "ssl-certificate"
	keySSLKey            = "ssl-key"
	keyAuth0Domain       = "auth0-domain"
	keyAuth0ClientID     = "auth0-client-id"
	keyAuth0ClientSecret = "auth0-client-secret"
	keyAuth0CallbackURL  = "auth0-callback-url"
)

// to be changed using ldflags with the go build command
var version = "unknown"

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	cmdAuth := &cobra.Command{
		Use: appName,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, varKey := range []string{
				keyAuth0Domain,
				keyAuth0ClientID,
				keyAuth0ClientSecret,
				keyAuth0CallbackURL,
			} {
				val := viper.GetString(varKey)
				if val == "" {
					return fmt.Errorf("%s not set", varKey)
				}
			}

			domain := viper.GetString(keyAuth0Domain)
			codeExchanger := &oauth2.Config{
				ClientID:     viper.GetString(keyAuth0ClientID),
				ClientSecret: viper.GetString(keyAuth0ClientSecret),
				RedirectURL:  viper.GetString(keyAuth0CallbackURL),
				Scopes:       []string{"openid", "email"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://" + domain + "/authorize",
					TokenURL: "https://" + domain + "/oauth/token",
				},
			}

			mux := monauth.ServeMux(codeExchanger, time.Second*10)
			handler := middleware.Logger(logger, mux, "")

			serveFn := newServeFn(
				logger,
				viper.GetString(keySSLCertificate),
				viper.GetString(keySSLKey),
			)
			addr := ":" + viper.GetString(keyPort)
			logger.Printf("Version: %s", version)
			logger.Printf("Serving at %s", addr)
			return serveFn(addr, handler)
		},
	}

	cmdAuth.AddCommand(versioncmd.New(version, os.Stdout))

	cobra.OnInitialize(viperAutoEnvVar)
	cmdAuth.Flags().String(keyPort, "80", "server listening port")
	cmdAuth.Flags().String(keySSLCertificate, "", "path to SSL certificate, leave empty for http")
	cmdAuth.Flags().String(keySSLKey, "", "path to SSL key, leave empty for http")
	cmdAuth.Flags().String(keyAuth0Domain, "", "auth0 domain to use for authentication")
	cmdAuth.Flags().String(keyAuth0ClientID, "", "auth0 client ID")
	cmdAuth.Flags().String(keyAuth0ClientSecret, "", "auth0 client secret")
	cmdAuth.Flags().String(keyAuth0CallbackURL, "", "auth0 callback URL")
	err := viper.BindPFlags(cmdAuth.Flags())
	if err != nil {
		logger.Printf("unable to BindPFlags: %v", err)
		os.Exit(1)
	}

	if err := cmdAuth.Execute(); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}

func viperAutoEnvVar() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
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
