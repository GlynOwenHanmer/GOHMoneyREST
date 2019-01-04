// +build functional

package functional

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/glynternet/mon/internal/client"
	"github.com/glynternet/mon/pkg/storage/postgres"
	"github.com/glynternet/mon/pkg/storage/storagetest"
)

const (
	keyServerHost = "SERVER_HOST"
	keyDBHost     = "DB_HOST"
	keyDBUser     = "DB_USER"
	keyDBPassword = "DB_PASSWORD"
	keyDBName     = "DB_NAME"
	keyDBSSLMode  = "DB_SSLMODE"
)

func TestMain(m *testing.M) {
	const retries = 10
	const backoff = 2 * time.Second
	errs := make([]error, retries)
	var i int
	for i = 0; i < retries; i++ {
		err := postgres.CreateStorage(
			os.Getenv(keyDBHost),
			os.Getenv(keyDBUser),
			os.Getenv(keyDBPassword),
			os.Getenv(keyDBName),
			os.Getenv(keyDBSSLMode),
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
	os.Exit(m.Run())
}

func TestSuite(t *testing.T) {
	host := os.Getenv(keyServerHost)
	store := client.Client(host)
	if !store.Available() {
		t.Fatalf("store at %q is unavailable", host)
	}
	storagetest.Test(t, store)
}
