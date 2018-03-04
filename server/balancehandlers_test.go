package server

import (
	"net/http"
	"testing"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting-storagetest"
	"github.com/glynternet/go-accounting/balance"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_balances(t *testing.T) {
	t.Run("SelectAccount error", func(t *testing.T) {
		expected := errors.New("account error")
		srv := &server{
			storage: &accountingtest.Storage{
				AccountErr: expected,
			},
		}
		code, bs, err := srv.balances(1) // any ID can be used because of the stub
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, bs)
	})

	t.Run("SelectBalance error", func(t *testing.T) {
		account := &storage.Account{}
		expected := errors.New("balances error")
		srv := &server{
			storage: &accountingtest.Storage{
				Account:     account,
				BalancesErr: expected,
			},
		}
		code, bs, err := srv.balances(1) // any ID can be used because of the stub
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Equal(t, expected, errors.Cause(err))
		assert.Nil(t, bs)
	})

	t.Run("all ok", func(t *testing.T) {
		expected := &storage.Balances{{ID: 1}}
		srv := &server{
			storage: &accountingtest.Storage{
				Account:  &storage.Account{},
				Balances: expected,
			},
		}
		code, bs, err := srv.balances(1) // any ID can be used because of the stub
		assert.Equal(t, http.StatusOK, code)
		assert.NoError(t, err)
		assert.IsType(t, new(storage.Balances), bs)
		assert.Equal(t, expected, bs)
	})
}

func TestServer_InsertBalance(t *testing.T) {
	t.Run("SelectAccount error", func(t *testing.T) {
		expected := errors.New("SelectAccount error")
		srv := server{&accountingtest.Storage{
			AccountErr: expected,
		}}
		code, b, err := srv.insertBalance(0, balance.Balance{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "selecting account")
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Nil(t, b)
	})

	t.Run("InsertBalance error", func(t *testing.T) {
		expected := errors.New("InsertBalance error")
		srv := server{&accountingtest.Storage{
			Account:    &storage.Account{},
			BalanceErr: expected,
		}}
		code, b, err := srv.insertBalance(0, balance.Balance{})
		assert.Equal(t, expected, errors.Cause(err))
		assert.Contains(t, err.Error(), "inserting balance")
		assert.Equal(t, http.StatusBadRequest, code)
		assert.Nil(t, b)
	})

	t.Run("all ok", func(t *testing.T) {
		expected := &storage.Balance{}
		srv := server{&accountingtest.Storage{
			Account: &storage.Account{},
			Balance: expected,
		}}
		code, b, err := srv.insertBalance(0, balance.Balance{})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, code)
		assert.Equal(t, expected, b)
	})
}
