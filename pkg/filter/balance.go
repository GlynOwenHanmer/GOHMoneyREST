package filter

import (
	"time"

	"github.com/glynternet/mon/pkg/storage"
)

// BalanceCondition is a function that will return true if a given storage.Balance
// satisfies some certain condition.
type BalanceCondition func(storage.Balance) bool

// After produces a BalanceCondition that can be used to identify if an
// is from after a given time.
func After(t time.Time) BalanceCondition {
	return func(a storage.Balance) bool {
		return a.Date.After(t)
	}
}
