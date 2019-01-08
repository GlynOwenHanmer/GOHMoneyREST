package client

import (
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/mon/internal/router"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/storagetest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestClient_SelectAccounts(t *testing.T) {
	s := &storagetest.Storage{
		Accounts: &storage.Accounts{
			{
				ID: 51,
				Account: *accountingtest.NewAccount(
					t,
					"test-0",
					accountingtest.NewCurrencyCode(t, "EUR"),
					time.Now().UTC().Truncate(time.Nanosecond),
				),
			},
			{
				ID: 981742,
				Account: *accountingtest.NewAccount(
					t,
					"test-1",
					accountingtest.NewCurrencyCode(t, "GBP"),
					// TODO: Revert this test to not have the UTC() and test when the timezone of the machine running is not UTC.
					time.Now().UTC().Add(time.Hour*123).Truncate(time.Nanosecond),
				),
			},
		},
	}

	router, listener, client := newTestComponents(t, s)

	errCh := make(chan error)
	go func() {
		errCh <- http.Serve(listener, router)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := client.SelectAccounts()
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Accounts, selected)
		close(errCh)
	}()

	common.FatalIfError(t, <-errCh, "received error")
}

func TestClient_SelectAccount(t *testing.T) {
	s := &storagetest.Storage{
		Account: &storage.Account{
			ID: 51,
			Account: *accountingtest.NewAccount(
				t,
				"test",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Now().UTC().Truncate(time.Nanosecond),
			),
		},
	}

	router, listener, client := newTestComponents(t, s)

	errCh := make(chan error)
	go func() {
		errCh <- http.Serve(listener, router)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := client.SelectAccount(734) // id doesn't matter when mocking
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Account, selected)
		close(errCh)
	}()

	common.FatalIfError(t, <-errCh, "received error")
}

func TestClient_SelectAccountBalances(t *testing.T) {
	s := &storagetest.Storage{
		Account: &storage.Account{
			ID: 51,
			Account: *accountingtest.NewAccount(
				t,
				"test",
				accountingtest.NewCurrencyCode(t, "EUR"),
				time.Now().Truncate(time.Nanosecond),
			),
		},
		Balances: &storage.Balances{
			storage.Balance{ID: 123},
		},
	}

	router, listener, client := newTestComponents(t, s)

	errCh := make(chan error)
	go func() {
		errCh <- http.Serve(listener, router)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		selected, err := client.SelectAccountBalances((*s).Account.ID) // id doesn't matter when mocking
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Equal(t, s.Balances, selected)
		close(errCh)
	}()

	common.FatalIfError(t, <-errCh, "received error")
}

func TestClient_InsertAccount(t *testing.T) {
	account := &storage.Account{
		ID: 51,
		Account: *accountingtest.NewAccount(
			t,
			"test",
			accountingtest.NewCurrencyCode(t, "EUR"),
			time.Date(3000, 0, 0, 0, 0, 0, 0, time.UTC),
		),
	}

	s := &storagetest.Storage{
		Account: account,
	}

	router, listener, client := newTestComponents(t, s)

	errCh := make(chan error)
	go func() {
		errCh <- http.Serve(listener, router)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		inserted, err := client.InsertAccount(account.Account)
		assert.NoError(t, err)
		assert.Equal(t, account.ID, inserted.ID)
		assert.Equal(t, account.Account, inserted.Account)
		close(errCh)
	}()

	common.FatalIfError(t, <-errCh, "received error")
}

func TestClient_InsertBalance(t *testing.T) {
	expected := &storage.Balance{ID: 293}
	now := time.Now()
	account := &storage.Account{
		Account: *accountingtest.NewAccount(t, "test account", accountingtest.NewCurrencyCode(t, "ABC"), now),
	}
	insert := balance.Balance{Date: now}

	s := &storagetest.Storage{
		Account: account,
		Balance: expected,
	}

	router, listener, client := newTestComponents(t, s)

	errCh := make(chan error)
	go func() {
		errCh <- http.Serve(listener, router)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		inserted, err := client.InsertBalance(0, insert, "test note")
		assert.NoError(t, err)
		assert.Equal(t, expected, inserted)
		assert.Equal(t, s.LastBalanceNote, "test note")
		close(errCh)
	}()

	common.FatalIfError(t, <-errCh, "received error")
}

func newTestComponents(t *testing.T, s storage.Storage) (*mux.Router, net.Listener, Client) {
	r := newTestRouter(t, s)
	l := newTestNetListener(t)
	c := newTestClient(l)
	return r, l, c
}

func newTestRouter(t *testing.T, s storage.Storage) *mux.Router {
	r, err := router.New(s, log.New(os.Stderr, "", log.LstdFlags))
	common.FatalIfError(t, err, "creating new router")
	if !assert.NotNil(t, r) {
		t.Fatal("expected non-nil router")
	}
	return r
}

func newTestNetListener(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", "localhost:0")
	common.FatalIfError(t, err, "creating new net listener")
	if !assert.NotNil(t, l) {
		t.Fatal("expected non-nil listener")
	}
	return l
}

func newTestClient(l net.Listener) Client {
	return Client("http://" + l.Addr().String())
}
