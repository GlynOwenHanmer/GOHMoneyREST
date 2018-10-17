package model

import "github.com/glynternet/mon/pkg/storage"

// SelectAccountBalances returns all Balances for a given Account and any
// errors that occur whilst attempting to retrieve the Balances. The Balances
// are sorted by chronological order then by the id of the Balance in the DB
func SelectAccountBalances(s storage.Storage, a storage.Account) (*storage.Balances, error) {
	return s.SelectAccountBalances(a.ID)
}
