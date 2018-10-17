package model_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/internal/model"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAccount(t *testing.T) {
	t.Run("select account balances error", func(t *testing.T) {
		expectedErr := errors.New("balances error")
		s := &storagetest.Storage{
			BalancesErr: expectedErr,
		}

		updated, actualErr := model.UpdateAccount(s, storage.Account{}, account.Account{})
		assert.Nil(t, updated)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, errors.Cause(actualErr))
		assert.Contains(t, actualErr.Error(), "selecting Account Balances for update validation")
	})

	now := time.Now()
	s := &storagetest.Storage{
		Balances: &storage.Balances{
			storage.Balance{
				Balance: balance.Balance{Date: now},
			},
		},
		AccountErr: errors.New("account error"),
	}

	a := accountingtest.NewAccount(t, "A", accountingtest.NewCurrencyCode(t, "YEN"), now)
	initial := storage.Account{
		ID:      999,
		Account: *a,
	}

	t.Run("invalid with balances", func(t *testing.T) {
		updates := accountingtest.NewAccount(t,
			"B",
			accountingtest.NewCurrencyCode(t, "GBP"),
			now.Add(time.Hour),
			account.CloseTime(now.Add(24*time.Hour)),
		)

		updated, err := model.UpdateAccount(s, initial, *updates)
		assert.Nil(t, updated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update would make balance invalid")
	})

	t.Run("update against storage", func(t *testing.T) {
		updates := accountingtest.NewAccount(t,
			"B",
			accountingtest.NewCurrencyCode(t, "GBP"),
			now.Add(-time.Hour),
			account.CloseTime(now.Add(24*time.Hour)),
		)

		_, err := model.UpdateAccount(s, initial, *updates)
		assert.Equal(t, initial.ID, s.LastAccountID)
		assert.Equal(t, s.AccountErr, errors.Cause(err))
	})
}
