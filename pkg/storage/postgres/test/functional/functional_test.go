// +build functional

package functional

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/glynternet/go-money/common"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/postgres"
	"github.com/glynternet/mon/pkg/storage/storagetest"
	"github.com/stretchr/testify/assert"
)

const (
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
			fmt.Printf("[retry: %02d] %v\n", i, err)
		}
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestSuite(t *testing.T) {
	store := createStorage(t)
	storagetest.Test(t, store)
}

func createStorage(t *testing.T) storage.Storage {
	cs, err := postgres.NewConnectionString(
		os.Getenv(keyDBHost),
		os.Getenv(keyDBUser),
		os.Getenv(keyDBPassword),
		os.Getenv(keyDBName),
		os.Getenv(keyDBSSLMode),
	)
	common.FatalIfError(t, err, "creating connection string")
	store, err := postgres.New(cs)
	common.FatalIfError(t, err, "creating storage")
	if !assert.True(t, store.Available(), "store should be available") {
		t.FailNow()
	}
	return store
}
