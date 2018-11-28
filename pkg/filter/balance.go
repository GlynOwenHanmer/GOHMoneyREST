package filter

import (
	"time"

	"github.com/glynternet/mon/pkg/storage"
)

// StorageBalanceCondition is a function that will return true if a given storage.Balance
// satisfies some certain condition.
type StorageBalanceCondition func(storage.Balance) bool

// Filter returns a set of storage.Balances that have been filtered down to the
// ones that match the StorageBalanceCondition
func (bc StorageBalanceCondition) Filter(bs storage.Balances) storage.Balances {
	var filtered storage.Balances
	for _, b := range bs {
		if bc(b) {
			filtered = append(filtered, b)
		}
	}
	return filtered
}

// StorageBalanceNot produces a StorageBalanceCondition that inverts the outcome of the given
// StorageBalanceCondition
func StorageBalanceNot(c StorageBalanceCondition) StorageBalanceCondition {
	return func(a storage.Balance) bool {
		return !c(a)
	}
}

// StorageBalanceAfter produces a StorageBalanceCondition that can be used to identify if a
// Balance is from after a given time.
func StorageBalanceAfter(t time.Time) StorageBalanceCondition {
	return func(a storage.Balance) bool {
		return a.Date.After(t)
	}
}
