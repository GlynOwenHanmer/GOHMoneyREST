package filter_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/pkg/filter"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func stubAccountCondition(match bool) filter.AccountCondition {
	return func(_ storage.Account) bool {
		return match
	}
}

func TestID(t *testing.T) {
	for _, test := range []struct {
		name string
		storage.Account
		id    uint
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "matching",
			Account: storage.Account{ID: 111},
			id:      111,
			match:   true,
		},
		{
			name:    "non-matching",
			Account: storage.Account{ID: 222},
			id:      123,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.ID(test.id)
			match := f(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestCurrency(t *testing.T) {

	for _, test := range []struct {
		name string
		storage.Account
		code  currency.Code
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "nil code with valid account",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "AUP"), time.Time{})},
		},
		{
			name: "nil account with valid code",
			code: accountingtest.NewCurrencyCode(t, "AUP"),
		},
		{
			name:    "valid code and account with non-matching code",
			code:    accountingtest.NewCurrencyCode(t, "BUP"),
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "AUP"), time.Time{})},
		},
		{
			name:    "valid code and account with matching code",
			code:    accountingtest.NewCurrencyCode(t, "BUP"),
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Time{})},
			match:   true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.Currency(test.code)
			match := f(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestExisted(t *testing.T) {
	for _, test := range []struct {
		name string
		storage.Account
		time.Time
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "time before open",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   false,
		},
		{
			name:    "time equal open",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   true,
		},
		{
			name:    "time after open",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(2001, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.Existed(test.Time)
			match := f(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestOpenAt(t *testing.T) {
	for _, test := range []struct {
		name string
		storage.Account
		time.Time
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name:    "open account",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC))},
			Time:    time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   true,
		},
		{
			name:    "closed account",
			Account: storage.Account{Account: *accountingtest.NewAccount(t, "test", accountingtest.NewCurrencyCode(t, "BUP"), time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC), account.CloseTime(time.Date(2001, 0, 0, 0, 0, 0, 0, time.UTC)))},
			Time:    time.Date(2002, 0, 0, 0, 0, 0, 0, time.UTC),
			match:   false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.OpenAt(test.Time)
			match := f(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestAccountConditions_Or(t *testing.T) {
	for _, test := range []struct {
		name string
		filter.AccountConditions
		storage.Account
		match bool
	}{
		{
			name: "zero-values",
		},
		{
			name:  "single condition matching",
			match: true,
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(true),
			},
		},
		{
			name: "single condition non-matching",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(false),
			},
		},
		{
			name:  "multiple conditions matching",
			match: true,
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(true),
				stubAccountCondition(true),
			},
		},
		{
			name: "multiple conditions non-matching",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(false),
				stubAccountCondition(false),
			},
		},
		{
			name:  "multiple conditions mixed-results",
			match: true,
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(false),
				stubAccountCondition(false),
				stubAccountCondition(true),
				stubAccountCondition(false),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			match := test.AccountConditions.Or(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}

func TestAccountConditions_And(t *testing.T) {
	for _, test := range []struct {
		name string
		filter.AccountConditions
		storage.Account
		match bool
	}{
		{
			name:  "zero-values",
			match: true,
		},
		{
			name: "single match",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(true),
			},
			match: true,
		},
		{
			name: "multiple match",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(true),
				stubAccountCondition(true),
			},
			match: true,
		},
		{
			name: "single non-match",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(false),
			},
		},
		{
			name: "multiple non-match",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(false),
				stubAccountCondition(false),
			},
		},
		{
			name: "starting with non-match and mixed others",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(false),
				stubAccountCondition(true),
				stubAccountCondition(false),
			},
		},
		{
			name: "starting with match and mixed others",
			AccountConditions: filter.AccountConditions{
				stubAccountCondition(true),
				stubAccountCondition(false),
				stubAccountCondition(true),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			match := test.AccountConditions.And(test.Account)
			assert.Equal(t, test.match, match)
		})
	}
}
func TestAccountCondition_Filter(t *testing.T) {
	matchingID := uint(1)
	nonmatchingID := uint(2)
	c := filter.ID(matchingID)
	for _, test := range []struct {
		name string
		in   storage.Accounts
		out  storage.Accounts
	}{
		{
			name: "zero-values",
		},
		{
			name: "single matching account",
			in:   storage.Accounts{{ID: matchingID}},
			out:  storage.Accounts{{ID: matchingID}},
		},
		{
			name: "single non-matching account",
			in:   storage.Accounts{{ID: nonmatchingID}},
		},
		{
			name: "single matching and single non-matching account",
			in:   storage.Accounts{{ID: nonmatchingID}, {ID: matchingID}},
			out:  storage.Accounts{{ID: matchingID}},
		},
		{
			name: "multiple mixed matching and non-matching accounts",
			in: storage.Accounts{
				{ID: nonmatchingID},
				{ID: nonmatchingID},
				{ID: matchingID},
				{ID: matchingID},
				{ID: nonmatchingID},
				{ID: matchingID},
				{ID: nonmatchingID},
			},
			out: storage.Accounts{
				{ID: matchingID},
				{ID: matchingID},
				{ID: matchingID},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			out := c.Filter(test.in)
			assert.Equal(t, test.out, out)
		})
	}
}
