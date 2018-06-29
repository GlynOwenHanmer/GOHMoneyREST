// Package storagetest provides some useful functionality for testing the storage package
package storagetest

import (
	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/storage"
)

// Storage is a data structure that satisfies the storage.Storage interface
type Storage struct {
	IsAvailable bool
	Err         error

	*storage.Account
	AccountErr error

	*storage.Accounts

	*storage.Balance
	BalanceErr error

	*storage.Balances
	BalancesErr error
}

// Available stubs storage.Available method
func (s *Storage) Available() bool { return s.IsAvailable }

// Close stubs the storage.Close method
func (s *Storage) Close() error    { return s.Err }

// InsertAccount stubs the storage.InsertAccount method
func (s *Storage) InsertAccount(account.Account) (*storage.Account, error) {
	return s.Account, s.AccountErr
}

// UpdateAccount stubs the storage.UpdateAccount method
func (s *Storage) UpdateAccount(a *storage.Account, updates *account.Account) (*storage.Account, error) {
	return s.Account, s.AccountErr
}

// SelectAccount stubs the storage.SelectAccount method
func (s *Storage) SelectAccount(uint) (*storage.Account, error) { return s.Account, s.AccountErr }
// SelectAccounts stubs the storage.SelectAccounts method
func (s *Storage) SelectAccounts() (*storage.Accounts, error)   { return s.Accounts, s.Err }

// DeleteAccount stubs the storage.DeleteAccount method
func (s *Storage) DeleteAccount(uint) error { return s.AccountErr }

// InsertBalance stubs the storage.InsertBalance method
func (s *Storage) InsertBalance(storage.Account, balance.Balance) (*storage.Balance, error) {
	return s.Balance, s.BalanceErr
}

// SelectAccountBalances stubs the storage.SelectAccountBalances method
func (s *Storage) SelectAccountBalances(storage.Account) (*storage.Balances, error) {
	return s.Balances, s.BalancesErr
}
