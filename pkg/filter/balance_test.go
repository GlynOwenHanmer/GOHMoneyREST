package filter_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/filter"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func stubBalanceCondition(match bool) filter.BalanceCondition {
	return func(_ storage.Balance) bool {
		return match
	}
}

func TestBalanceNot(t *testing.T) {
	var dummy storage.Balance
	for _, test := range []struct {
		name string
		in   bool
	}{
		{
			name: "zero-values",
		},
		{
			name: "inner filter match",
			in:   true,
		},
		{
			name: "inner filter non-match",
			in:   false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := filter.BalanceNot(stubBalanceCondition(test.in))
			match := f(dummy)
			assert.Equal(t, !test.in, match)
		})
	}
}

func TestAfter(t *testing.T) {
	b := storage.Balance{
		Balance: balance.Balance{
			Date: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}

	tests := []struct {
		name string
		time.Time
		match bool
	}{
		{
			name:  "after",
			match: true,
			Time:  time.Date(1000, 1, 1, 1, 1, 1, 1, time.UTC),
		},
		{
			name: "not after",
			Time: time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := filter.After(tt.Time)(b)
			assert.Equal(t, tt.match, match)
		})
	}
}

func TestBalanceCondition_Filter(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	c := filter.After(date)
	for _, test := range []struct {
		name string
		in   storage.Balances
		out  storage.Balances
	}{
		{
			name: "zero-values",
		},
		{
			name: "single matching balance",
			in:   storage.Balances{{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}}},
			out:  storage.Balances{{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}}},
		},
		{
			name: "single non-matching balance",
			in:   storage.Balances{{Balance: balance.Balance{Date: date.Add(-100 * time.Hour)}}},
		},
		{
			name: "single matching and single non-matching balance",
			in: storage.Balances{
				{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(-100 * time.Hour)}},
			},
			out: storage.Balances{{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}}},
		},
		{
			name: "multiple mixed matching and non-matching accounts",
			in: storage.Balances{
				{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(200 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(-100 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(-200 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(170 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(-190 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(1200 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(-100 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(-100 * time.Hour)}},
			},
			out: storage.Balances{
				{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(200 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(170 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(1200 * time.Hour)}},
				{Balance: balance.Balance{Date: date.Add(100 * time.Hour)}},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			out := c.Filter(test.in)
			assert.Equal(t, test.out, out)
		})
	}
}
