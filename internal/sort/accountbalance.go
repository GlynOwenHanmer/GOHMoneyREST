package sort

import (
	"sort"

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

// BalanceAmount sorts a slice of accountbalance.AccountBalance by the amount
// of the Balance, in ascending order.
// BalanceAmount cannot guarantee any specific order within a subsection of
// the slice when multiple AccountBalance have the same amount.
func BalanceAmount(abs []accountbalance.AccountBalance) {
	sortAccountBalances(abs, balanceAmount)
}

func balanceAmount(a, b accountbalance.AccountBalance) bool {
	return a.Balance.Amount < b.Balance.Amount
}

// BalanceAmountMagnitude sorts a slice of accountbalance.AccountBalance by the
// absolute magnitude of the amount of the Balance, in ascending order.
// BalanceAmountMagnitude cannot guarantee any specific order within a
// subsection of the slice when multiple AccountBalance have the same absolute
// amount.
func BalanceAmountMagnitude(abs []accountbalance.AccountBalance) {
	sortAccountBalances(abs, balanceMagnitude)
}

func balanceMagnitude(a, b accountbalance.AccountBalance) bool {
	absI := absolute(a.Balance.Amount)
	absJ := absolute(b.Balance.Amount)
	return absI < absJ
}

func absolute(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
