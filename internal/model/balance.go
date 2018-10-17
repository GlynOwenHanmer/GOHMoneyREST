package model

import (
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
)

// SelectAccountBalances returns all Balances for a given Account and any
// errors that occur whilst attempting to retrieve the Balances. The Balances
// are sorted by chronological order then by the id of the Balance in the DB
func SelectAccountBalances(s storage.Storage, a storage.Account) (*storage.Balances, error) {
	return s.SelectAccountBalances(a.ID)
}

// InsertBalance will insert a Balance into the given storage using the value
// of the given storage.Account. InsertBalance will perform any logic checks
// before attempting to insert the balance into the given Storage.
func InsertBalance(s storage.Storage, a storage.Account, b balance.Balance) (*storage.Balance, error) {
	err := a.Account.ValidateBalance(b)
	if err != nil {
		return nil, errors.Wrap(err, "validating balance")
	}
	dbb, err := s.InsertBalance(a.ID, b)
	return dbb, errors.Wrap(err, "inserting balance")
}
