package storage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/stretchr/testify/assert"
)

func TestBalance_Equal(t *testing.T) {
	for _, test := range []struct {
		name  string
		a, b  Balance
		equal bool
	}{
		{
			name:  "zero-values",
			equal: true,
		},
		{
			name: "unequal IDs",
			a:    Balance{ID: 1},
			b:    Balance{ID: 2},
		},
		{
			name: "unequal inner Balance",
			a:    Balance{Balance: balance.Balance{Amount: 1}},
			b:    Balance{Balance: balance.Balance{Amount: 2}},
		},
		{
			name: "unequal Note",
			a:    Balance{ID: 4, Balance: balance.Balance{Amount: 1}, Note: "this note"},
			b:    Balance{ID: 4, Balance: balance.Balance{Amount: 1}},
		},
		{
			name:  "equal",
			a:     Balance{ID: 4, Balance: balance.Balance{Amount: 1}},
			b:     Balance{ID: 4, Balance: balance.Balance{Amount: 1}},
			equal: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.equal, test.a.Equal(test.b))
		})
	}
}

func TestBalance_JSONLoop(t *testing.T) {
	b := Balance{
		ID:   47,
		Note: "test note",
		Balance: balance.Balance{
			Date:   time.Date(1000, 0, 0, 0, 0, 0, 0, time.UTC),
			Amount: -34567,
		},
	}
	bs, err := json.Marshal(b)
	common.FatalIfError(t, err, "marshalling Balance json")
	var actual Balance
	err = json.Unmarshal(bs, &actual)
	assert.Nil(t, err)
	assert.Equal(t, b, actual)
}

func TestBalances_InnerBalances(t *testing.T) {
	for _, test := range []struct {
		name string
		in   Balances
		out  balance.Balances
	}{
		{
			name: "zero-values",
		},
		{
			name: "single balance",
			in:   Balances{{}},
			out:  balance.Balances{{}},
		},
		{
			name: "multiple balance",
			in: Balances{
				{
					Balance: balance.Balance{
						Date: time.Date(2, 2, 2, 2, 2, 2, 2, time.UTC)},
				},
				{
					Balance: balance.Balance{
						Date: time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)},
				},
			},
			out: balance.Balances{
				balance.Balance{
					Date: time.Date(2, 2, 2, 2, 2, 2, 2, time.UTC),
				},
				balance.Balance{
					Date: time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			out := test.in.InnerBalances()
			assert.Equal(t, test.out, out)
		})
	}
}
