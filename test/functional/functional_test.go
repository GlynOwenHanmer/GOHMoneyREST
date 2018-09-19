// +build functional

package functional

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/glynternet/mon/internal/client"
	"github.com/glynternet/mon/pkg/storage/postgres"
	"github.com/glynternet/mon/pkg/storage/storagetest"
	"github.com/spf13/viper"
)

const (
	// viper keys
	keyServerHost = "server-host"
	keyDBHost     = "db-host"
	keyDBUser     = "db-user"
	keyDBName     = "db-name"
	keyDBSSLMode  = "db-sslmode"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
	setup()
}

func setup() {
	const retries = 10
	const backoff = 2 * time.Second
	errs := make([]error, retries)
	var i int
	for i = 0; i < retries; i++ {
		err := postgres.CreateStorage(
			viper.GetString(keyDBHost),
			viper.GetString(keyDBUser),
			viper.GetString(keyDBName),
			viper.GetString(keyDBSSLMode),
		)
		if err == nil {
			break
		}
		errs[i] = err
		time.Sleep(backoff)
	}
	if errs[retries-1] != nil {
		for i, err := range errs {
			log.Printf("[retry: %02d] %v\n", i, err)
		}
		os.Exit(1)
	}
	log.Print("Setup complete")
}

func TestSuite(t *testing.T) {
	host := viper.GetString(keyServerHost)
	store := client.Client(host)
	if !store.Available() {
		t.Fatalf("store at %q is unavailable", host)
	}
	storagetest.Test(t, store)
}
