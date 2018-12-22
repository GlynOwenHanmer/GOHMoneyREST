package storagetest

import (
	"strconv"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const numOfAccounts = 2

// Test will run a suite of tests again a given Storage
func Test(t *testing.T, store storage.Storage) {
	tests := []struct {
		title string
		run   func(t *testing.T, c storage.Storage)
	}{
		{
			title: "inserting and retrieving accounts",
			run:   insertAndRetrieveAccounts,
		},
		{
			title: "inserting and retrieving balances",
			run:   insertDeleteAndRetrieveBalances,
		},
		{
			title: "update account",
			run:   updateAccount,
		},
		{
			title: "insert and delete accounts",
			run:   insertAndDeleteAccounts,
		},
	}
	for _, test := range tests {
		success := t.Run(test.title, func(t *testing.T) {
			test.run(t, store)
		})
		if !success {
			t.Fail()
			return
		}
	}
}

func insertAndRetrieveAccounts(t *testing.T, store storage.Storage) {
	as, err := store.SelectAccounts()
	common.FatalIfError(t, err, "selecting accounts")

	if !assert.Len(t, *as, 0) {
		t.FailNow()
	}

	a := accountingtest.NewAccount(t, "A", accountingtest.NewCurrencyCode(t, "YEN"), time.Now())
	insertedA, err := store.InsertAccount(*a)
	common.FatalIfError(t, err, "inserting account")

	as = selectAccounts(t, store)

	if !assert.Len(t, *as, 1) {
		t.FailNow()
	}
	selectedA := &(*as)[0]
	equal, err := insertedA.Equal(*selectedA)
	common.FatalIfError(t, err, "equaling inserted and retrieved")
	if !assert.True(t, equal) {
		t.FailNow()
	}

	selectedByIDA, err := store.SelectAccount(insertedA.ID)
	common.FatalIfError(t, err, "selecting account by ID")

	assertThreeAccountsEqual(t, insertedA, selectedA, selectedByIDA)

	b := accountingtest.NewAccount(t,
		"B",
		accountingtest.NewCurrencyCode(t, "EUR"),
		time.Now().Add(-1*time.Hour),
	)

	insertedB, err := store.InsertAccount(*b)
	common.FatalIfError(t, err, "inserting account")

	as, err = store.SelectAccounts()
	common.FatalIfError(t, err, "selecting accounts after inserting two")

	if !assert.Len(t, *as, numOfAccounts) {
		t.FailNow()
	}
	selectedB := &(*as)[1]
	equal, err = insertedB.Equal(*selectedB)
	common.FatalIfError(t, err, "equaling inserted and retrieved")
	if !assert.True(t, equal) {
		t.FailNow()
	}

	selectedByIDB, err := store.SelectAccount(insertedB.ID)
	common.FatalIfError(t, err, "selecting account by ID")

	assertThreeAccountsEqual(t, insertedB, selectedB, selectedByIDB)

	equal, err = insertedA.Equal(*insertedB)
	common.FatalIfError(t, err, "equaling insertedA and insertedB")
	if !assert.False(t, equal) {
		t.FailNow()
	}

	equal, err = selectedA.Equal(*selectedB)
	common.FatalIfError(t, err, "equaling selectedA and selectedB")
	if !assert.False(t, equal) {
		t.FailNow()
	}
}

// assertThreeAccountsEqual will compare three accounts against each other and fail the
// test if any of them are not equal against each other.
func assertThreeAccountsEqual(t *testing.T, a, b, c *storage.Account) {
	for _, as := range []struct {
		A, B *storage.Account
	}{
		{
			A: a, B: b,
		},
		{
			A: a, B: c,
		},
		{
			A: b, B: c,
		},
	} {
		equal, err := as.A.Equal(*as.B)
		common.FatalIfErrorf(t, err, "equalling accounts %+v", as)
		if !equal {
			t.Fatal("not equal")
		}
	}
}

func insertDeleteAndRetrieveBalances(t *testing.T, store storage.Storage) {
	as := selectAccounts(t, store)
	assert.Len(t, *as, numOfAccounts)
	for _, a := range *as {
		// assert that all accounts contain no balances
		bs, err := store.SelectAccountBalances(a.ID)
		common.FatalIfError(t, err, "selecting account balances")
		assert.Len(t, *bs, 0)

		// insert single balance
		// TODO: work out what is best to do here. Unfortunately, due to how
		// TODO: fine grained that postgres can store times, we cannot just use
		// TODO: time.Now() or some other fine grained time, as the balance
		// TODO: that is selected will not have the same time.
		// By using the a.Account.Opened() here, we know that the time is supported by the storage.
		b := newTestBalance(t, a.Account.Opened())
		note := "first balance"
		inserted, err := store.InsertBalance(a.ID, b, note)
		common.FatalIfError(t, err, "inserting Balance")
		equal := b.Equal(inserted.Balance)
		if !assert.True(t, equal) {
			t.FailNow()
		}
		if !assert.Equal(t, note, inserted.Note) {
			t.FailNow()
		}

		bs, err = store.SelectAccountBalances(a.ID)
		common.FatalIfError(t, err, "selecting account balances")
		assert.Len(t, *bs, 1)

		// delete balance
		err = store.DeleteBalance(inserted.ID)
		assert.NoError(t, err)

		bs, err = store.SelectAccountBalances(a.ID)
		common.FatalIfError(t, err, "selecting account balances")
		assert.Len(t, *bs, 0)
	}
}

func updateAccount(t *testing.T, store storage.Storage) {
	initial := accountingtest.NewAccount(t, "A", accountingtest.NewCurrencyCode(t, "YEN"), time.Now())

	inserted, err := store.InsertAccount(*initial)
	common.FatalIfError(t, err, "inserting account to store")

	// Here we truncate to the closest second to avoid the issue where
	// postgres stores times down to only the closest millisecond or so
	// TODO: Sort out rounding of times logic and document it properly. For
	// TODO: the moment, it is assumed that accounts won't be updated and
	// TODO: then compared against their original down to such a fine grain
	updates := accountingtest.NewAccount(t,
		"B",
		accountingtest.NewCurrencyCode(t, "GBP"),
		time.Now().Truncate(time.Second),
		account.CloseTime(time.Now().Add(24*time.Hour).Truncate(time.Second)),
	)

	updatedA, err := store.UpdateAccount(inserted.ID, *updates)
	common.FatalIfError(t, err, "updating account")
	if !assert.NotNil(t, updatedA) {
		t.FailNow()
	}
	assert.Equal(t, updatedA.ID, inserted.ID)
	assert.True(t,
		updatedA.Account.Equal(*updates),
		"inserted: %+v\nupdates: %+v\nupdatedA: %+v", inserted.Account, updates, updatedA.Account,
	)
}

func insertAndDeleteAccounts(t *testing.T, store storage.Storage) {
	selectedBefore := selectAccounts(t, store)

	type AccountBalances struct {
		storage.Account
		storage.Balances
	}
	var abs []AccountBalances
	for _, a := range *selectedBefore {
		bs, err := store.SelectAccountBalances(a.ID)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "selecting account balances for account %+v", a))
		}
		abs = append(abs, AccountBalances{
			Account:  a,
			Balances: *bs,
		})
	}

	const numInserted = 5

	var as []storage.Account
	for i := 0; i < numInserted; i++ {
		a := accountingtest.NewAccount(t, "TO DELETE", accountingtest.NewCurrencyCode(t, "BBC"), time.Now())
		ia, err := store.InsertAccount(*a)
		common.FatalIfError(t, err, "inserting account")
		as = append(as, *ia)
	}

	selectedAfter := selectAccounts(t, store)
	assert.Len(t, *selectedAfter, len(*selectedBefore)+numInserted)

	for i, a := range as {
		t.Run("deleting account (i:"+strconv.Itoa(i)+") should reduce accounts count by 1", func(t *testing.T) {
			err := store.DeleteAccount(a.ID)
			selectedAfter = selectAccounts(t, store)
			common.FatalIfError(t, err, "deleting account")
			// Accounts count should be the number of originals, with the
			// number that were inserted, then -1 for every delete
			assert.Len(t, *selectedAfter, len(*selectedBefore)+numInserted-(i+1))
		})
	}

	t.Run("previously existing accounts should remain the same", func(t *testing.T) {
		if !assert.Len(t, *selectedAfter, len(abs)) {
			t.FailNow()
		}
		for i := range *selectedAfter {
			afterAccount := (*selectedAfter)[i]
			assert.Equal(t, afterAccount, abs[i].Account)
			afters, err := store.SelectAccountBalances(afterAccount.ID)
			if err != nil {
				t.Fatal(errors.Wrapf(err, "selecting account balances for account %+v", afterAccount))
			}
			assert.Equal(t, *afters, abs[i].Balances)
		}
	})
}

func selectAccounts(t *testing.T, store storage.Storage) *storage.Accounts {
	as, err := store.SelectAccounts()
	common.FatalIfError(t, err, "selecting accounts after inserting one")
	return as
}

func newTestBalance(t *testing.T, time time.Time, os ...balance.Option) balance.Balance {
	b, err := balance.New(time, os...)
	common.FatalIfError(t, err, "creating test balance")
	return *b
}
