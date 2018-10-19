package sort

import (
	"sort"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/internal/accountbalance"
)

// sortAccountBalances sorts a slice of accountbalance.AccountBalance into the
// order determined by the given accountBalanceComparison, c.
func sortAccountBalances(abs []accountbalance.AccountBalance, c accountBalanceComparison) {
	sort.Slice(abs, func(i, j int) bool {
		return c(abs[i], abs[j])
	})
}

// accountBalanceComparison should return true if AccountBalance a should be
// before AccountBalance b in the sorted order.
type accountBalanceComparison func(a, b accountbalance.AccountBalance) bool

// balanceComparison should return true if Balance a should be before Balance b
// in the sorted order.
type balanceComparison func(a, b balance.Balance) bool

// newAccountBalanceComparison creates an accountBalanceComparison from a given
// balanceComparison. The generated accountBalanceComparison would provide the
// order purely based on the result of the given balanceComparison.
func newAccountBalanceComparison(c balanceComparison) accountBalanceComparison {
	return func(a, b accountbalance.AccountBalance) bool {
		return c(a.Balance, b.Balance)
	}
}

// BalanceAmount sorts a slice of accountbalance.AccountBalance by the amount
// of the Balance, in ascending order.
// BalanceAmount cannot guarantee any specific order within a subsection of
// the slice when multiple AccountBalance have the same amount.
func BalanceAmount(abs []accountbalance.AccountBalance) {
	sortAccountBalances(abs, newAccountBalanceComparison(balanceAmount))
}

func balanceAmount(a, b balance.Balance) bool {
	return a.Amount < b.Amount
}

// BalanceAmountMagnitude sorts a slice of accountbalance.AccountBalance by the
// absolute magnitude of the amount of the Balance, in ascending order.
// BalanceAmountMagnitude cannot guarantee any specific order within a
// subsection of the slice when multiple AccountBalance have the same absolute
// amount.
func BalanceAmountMagnitude(abs []accountbalance.AccountBalance) {
	sortAccountBalances(abs, newAccountBalanceComparison(balanceMagnitude))
}

func balanceMagnitude(a, b balance.Balance) bool {
	absA := absolute(a.Amount)
	absB := absolute(b.Amount)
	return absA < absB
}

func absolute(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
