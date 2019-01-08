package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/glynternet/mon/internal/router"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "monserve"

	// viper keys
	keyPort           = "port"
	keySSLCertificate = "ssl-certificate"
	keySSLKey         = "ssl-key"
	keyDBHost         = "db-host"
	keyDBUser         = "db-user"
	keyDBPassword     = "db-password"
	keyDBName         = "db-name"
	keyDBSSLMode      = "db-sslmode"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	var cmdDBServe = &cobra.Command{
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
			r, err := router.New(store, logger)
			if err != nil {
				return errors.Wrap(err, "error creating new server")
			}

			serveFn := newServeFn(
				logger,
				viper.GetString(keySSLCertificate),
				viper.GetString(keySSLKey),
			)
			addr := ":" + viper.GetString(keyPort)
			logger.Printf("Serving at %s", addr)
			return serveFn(addr, r)
		},
	}

	cobra.OnInitialize(viperAutoEnvVar)
	cmdDBServe.Flags().String(keyPort, "80", "server listening port")
	cmdDBServe.Flags().String(keySSLCertificate, "", "path to SSL certificate, leave empty for http")
	cmdDBServe.Flags().String(keySSLKey, "", "path to SSL key, leave empty for https")
	cmdDBServe.Flags().String(keyDBHost, "", "host address of the DB backend")
	cmdDBServe.Flags().String(keyDBName, "", "name of the DB set to use")
	cmdDBServe.Flags().String(keyDBUser, "", "DB user to authenticate with")
	cmdDBServe.Flags().String(keyDBPassword, "", "DB password to authenticate with")
	cmdDBServe.Flags().String(keyDBSSLMode, "", "DB SSL mode to use")
	err := viper.BindPFlags(cmdDBServe.Flags())
	if err != nil {
		logger.Printf("unable to BindPFlags: %v", err)
		os.Exit(1)
	}

	if err := cmdDBServe.Execute(); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
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
