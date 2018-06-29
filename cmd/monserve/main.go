package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/postgres"
	"github.com/glynternet/mon/router"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "monserve"

	// viper keys
	keyPort      = "port"
	keyDBHost    = "db-host"
	keyDBUser    = "db-user"
	keyDBName    = "db-name"
	keyDBSSLMode = "db-sslmode"
)

func main() {
	log.Fatal(cmdDBServe.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cmdDBServe.Flags().String(keyPort, "80", "server listening port")
	cmdDBServe.Flags().String(keyDBHost, "", "host address of the DB backend")
	cmdDBServe.Flags().String(keyDBName, "", "name of the DB set to use")
	cmdDBServe.Flags().String(keyDBUser, "", "DB user to authenticate with")
	cmdDBServe.Flags().String(keyDBSSLMode, "", "DB SSL mode to use")
	err := viper.BindPFlags(cmdDBServe.Flags())
	if err != nil {
		log.Printf("unable to BindPFlags: %v", err)
	}
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func newStorage(host, user, dbname, sslmode string) (storage.Storage, error) {
	cs, err := postgres.NewConnectionString(host, user, dbname, sslmode)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection string: %v", err)
	}
	return postgres.New(cs)
}

var cmdDBServe = &cobra.Command{
	Use: appName,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := newStorage(
			viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode),
		)
		if err != nil {
			return errors.Wrap(err, "error creating storage")
		}
		r, err := router.New(store)
		if err != nil {
			return errors.Wrap(err, "error creating new server")
		}
		return http.ListenAndServe(":"+viper.GetString(keyPort), r)
	},
}
