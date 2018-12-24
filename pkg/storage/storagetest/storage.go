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

	LastAccountID   uint
	LastBalanceNote string
}

// Available stubs storage.Available method
func (s *Storage) Available() bool { return s.IsAvailable }

// Close stubs the storage.Close method
func (s *Storage) Close() error { return s.Err }

// InsertAccount stubs the storage.InsertAccount method
func (s *Storage) InsertAccount(account.Account) (*storage.Account, error) {
	return s.Account, s.AccountErr
}

// UpdateAccount stubs the storage.UpdateAccount method
func (s *Storage) UpdateAccount(id uint, updates account.Account) (*storage.Account, error) {
	s.LastAccountID = id
	return s.Account, s.AccountErr
}

// SelectAccount stubs the storage.SelectAccount method
func (s *Storage) SelectAccount(id uint) (*storage.Account, error) {
	s.LastAccountID = id
	return s.Account, s.AccountErr
}

// SelectAccounts stubs the storage.SelectAccounts method
func (s *Storage) SelectAccounts() (*storage.Accounts, error) { return s.Accounts, s.Err }

// DeleteAccount stubs the storage.DeleteAccount method
func (s *Storage) DeleteAccount(id uint) error {
	s.LastAccountID = id
	return s.AccountErr
}

// InsertBalance stubs the storage.InsertBalance method
func (s *Storage) InsertBalance(accountID uint, _ balance.Balance, note string) (*storage.Balance, error) {
	s.LastAccountID = accountID
	s.LastBalanceNote = note
	return s.Balance, s.BalanceErr
}

// DeleteBalance stubs the storage.DeleteBalance method
func (s *Storage) DeleteBalance(_ uint) error { return s.Err }

// SelectAccountBalances mocks the storage.SelectAccountBalances method
func (s *Storage) SelectAccountBalances(id uint) (*storage.Balances, error) {
	s.LastAccountID = id
	return s.Balances, s.BalancesErr
}
