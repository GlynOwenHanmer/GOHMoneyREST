package sort

import (
	"sort"

	"github.com/glynternet/mon/pkg/storage"
)

// sortAccounts sorts a slice of storage.Accounts into the order determined by
// the given accountComparison, c.
func sortAccounts(as storage.Accounts, c accountComparison) {
	sort.Slice(as, func(i, j int) bool {
		return c(as[i], as[j])
	})
}

// accountComparison should return true if Account a should be before Account
// b in the sorted order.
type accountComparison func(a, b storage.Account) bool

// AccountID sorts a storage.Accounts by AccountID in ascending order.
// AccountID cannot guarantee any specific order within a subsection of
// storage.Accounts when multiple accounts have the same AccountID.
func AccountID(as storage.Accounts) {
	sortAccounts(as, accountID)
}

func accountID(a, b storage.Account) bool {
	return a.ID < b.ID
}

// AccountName sorts a storage.Accounts by AccountName in ascending order.
// AccountName cannot guarantee any specific order within a subsection of
// storage.Accounts when multiple accounts have the same name.
func AccountName(as storage.Accounts) {
	sortAccounts(as, accountName)
}

func accountName(a, b storage.Account) bool {
	return a.Account.Name() < b.Account.Name()
}
