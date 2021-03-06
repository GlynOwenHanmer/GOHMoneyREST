package storage

import (
	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
)

// Storage is something that can be used to store certain go-accounting types
type Storage interface {
	Available() bool
	Close() error
	InsertAccount(a account.Account) (*Account, error)
	SelectAccount(id uint) (*Account, error)
	UpdateAccount(id uint, updates account.Account) (*Account, error)
	SelectAccounts() (*Accounts, error)
	DeleteAccount(id uint) error
	//
	InsertBalance(accountID uint, b balance.Balance, note string) (*Balance, error)
	SelectAccountBalances(id uint) (*Balances, error)
	//UpdateBalance(a Account, b *Balance, us balance.Balance) error
	DeleteBalance(id uint) error
}
