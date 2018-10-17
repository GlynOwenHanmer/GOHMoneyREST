package filter

import (
	"time"

	"github.com/glynternet/mon/pkg/storage"
)

// BalanceCondition is a function that will return true if a given storage.Balance
// satisfies some certain condition.
type BalanceCondition func(storage.Balance) bool

// Filter returns a set of storage.Balances that have been filtered down to the
// ones that match the BalanceCondition
func (bc BalanceCondition) Filter(bs storage.Balances) storage.Balances {
	var filtered storage.Balances
	for _, b := range bs {
		if bc(b) {
			filtered = append(filtered, b)
		}
	}
	return filtered
}

// BalanceNot produces a BalanceCondition that inverts the outcome of the given
// BalanceCondition
func BalanceNot(c BalanceCondition) BalanceCondition {
	return func(a storage.Balance) bool {
		return !c(a)
	}
}

// After produces a BalanceCondition that can be used to identify if an
// is from after a given time.
func After(t time.Time) BalanceCondition {
	return func(a storage.Balance) bool {
		return a.Date.After(t)
	}
}
