package filter

import (
	"time"

	"github.com/glynternet/go-accounting/balance"
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

// BalanceCondition is a function that will return true if a given balance.Balance
// satisfies some certain condition.
type BalanceCondition func(balance.Balance) bool

// Filter returns a set of balance.Balances that have been filtered down to the
// ones that match the BalanceCondition
func (bc BalanceCondition) Filter(bs balance.Balances) balance.Balances {
	var filtered balance.Balances
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
	return func(b balance.Balance) bool {
		return !c(b)
	}
}

// BalanceAfter produces a BalanceCondition that can be used to identify if a
// Balance is from after a given time.
func BalanceAfter(t time.Time) BalanceCondition {
	return func(a balance.Balance) bool {
		return a.Date.After(t)
	}
}
