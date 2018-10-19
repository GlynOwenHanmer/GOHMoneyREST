package sort_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/mon/internal/sort"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	t.Run("nil input remains nil", func(t *testing.T) {
		var as storage.Accounts
		sort.AccountID(as)
		assert.Nil(t, as)
	})

	for _, test := range []struct {
		name    string
		in, out []uint
	}{
		{
			name: "empty storage.Accounts",
			in:   []uint{},
		},
		{
			name: "single account yields same account",
			in:   []uint{1},
			out:  []uint{1},
		},
		{
			name: "in-order input yields same order output",
			in:   []uint{0, 1, 2, 3},
			out:  []uint{0, 1, 2, 3},
		},
		{
			name: "in reverse order",
			in:   []uint{3, 2, 1, 0},
			out:  []uint{0, 1, 2, 3},
		},
		{
			name: "mixed order",
			in:   []uint{3, 7, 99, 2, 5, 8},
			out:  []uint{2, 3, 5, 7, 8, 99},
		},
		{
			name: "mixed order with repeating",
			in:   []uint{9, 5, 2, 66, 4, 5},
			out:  []uint{2, 4, 5, 5, 9, 66},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var as storage.Accounts
			for _, id := range test.in {
				as = append(as, storage.Account{ID: id})
			}

			sort.AccountID(as)

			var out []uint
			for _, a := range as {
				out = append(out, a.ID)
			}

			assert.Equal(t, test.out, out)
		})
	}
}

func TestName(t *testing.T) {
	t.Run("nil input remains nil", func(t *testing.T) {
		var as storage.Accounts
		sort.AccountName(as)
		assert.Nil(t, as)
	})

	namedAccount := func(name string) storage.Account {
		a, err := account.New(name, nil, time.Time{})
		common.FatalIfError(t, err, "creating new account")
		return storage.Account{Account: *a}
	}

	for _, test := range []struct {
		name    string
		in, out []string
	}{
		{
			name: "empty storage.Accounts",
			in:   []string{},
		},
		{
			name: "single account yields same account",
			in:   []string{"A"},
			out:  []string{"A"},
		},
		{
			name: "in-order input yields same order output",
			in:   []string{"A", "B", "C", "D"},
			out:  []string{"A", "B", "C", "D"},
		},
		{
			name: "in reverse order",
			in:   []string{"D", "C", "B", "A"},
			out:  []string{"A", "B", "C", "D"},
		},
		{
			name: "mixed order",
			in:   []string{"D", "C", "A", "B", "a"},
			out:  []string{"A", "B", "C", "D", "a"},
		},
		{
			name: "mixed order with repeating",
			in:   []string{"D", "C", "A", "B", "a", "D"},
			out:  []string{"A", "B", "C", "D", "D", "a"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var as storage.Accounts
			for _, name := range test.in {
				as = append(as, namedAccount(name))
			}

			sort.AccountName(as)

			var out []string
			for _, a := range as {
				out = append(out, a.Account.Name())
			}

			assert.Equal(t, test.out, out)
		})
	}
}
