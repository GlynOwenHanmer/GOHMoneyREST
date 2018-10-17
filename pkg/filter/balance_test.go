package filter

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/stretchr/testify/assert"
)

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
			match := After(tt.Time)(b)
			assert.Equal(t, tt.match, match)
		})
	}
}
