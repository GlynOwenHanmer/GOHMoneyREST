package model_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/accountingtest"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/internal/model"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/storage/storagetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInsertBalance(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		b, err := model.InsertBalance(nil, storage.Account{}, balance.Balance{})
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validating balance")
	})

	t.Run("id is passed", func(t *testing.T) {
		s := &storagetest.Storage{
			Balance:    &storage.Balance{},
			BalanceErr: errors.New("insert balance error"),
		}
		now := time.Now()
		b, err := model.InsertBalance(s,
			storage.Account{
				ID: 9183,
				Account: *accountingtest.NewAccount(t,
					"test account",
					accountingtest.NewCurrencyCode(t, "ABC"),
					now),
			},
			balance.Balance{Date: now})
		assert.Equal(t, s.BalanceErr, errors.Cause(err))
		assert.Equal(t, s.Balance, b)
		assert.Equal(t, uint(9183), s.LastAccountID)
	})
}
