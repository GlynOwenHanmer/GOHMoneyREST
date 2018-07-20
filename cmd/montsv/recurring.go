package main

import (
	"fmt"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	"github.com/pkg/errors"
)

const maxMonthlyDate = 28

type balanceGenerator interface {
	generateBalance(time time.Time) (*balance.Balance, error)
	generateAccountBalances(times []time.Time) (AccountBalances, error)
}

type dailyRecurringCost struct {
	name string
	currency.Code
	Amount int
	from   time.Time
}

func (rcs dailyRecurringCost) generateBalance(at time.Time) (*balance.Balance, error) {
	var amount int
	if at.After(rcs.from) {
		amount = int(at.Sub(rcs.from)/(time.Hour*24)) * rcs.Amount
	}
	b, err := balance.New(at, balance.Amount(amount))
	return b, errors.Wrap(err, "creating balance")
}

func (rcs dailyRecurringCost) generateAccountBalances(times []time.Time) (AccountBalances, error) {
	a, err := account.New(rcs.name, rcs.Code, time.Time{}) // time/date of account is not used currently
	if err != nil {
		return AccountBalances{}, errors.Wrap(err, "creating new account")
	}
	var bs balance.Balances
	for _, t := range times {
		b, err := rcs.generateBalance(t)
		if err != nil {
			return AccountBalances{}, errors.Wrapf(err, "generating balance for time:%s", t)
		}
		bs = append(bs, *b)
	}
	return AccountBalances{
		Account:  *a,
		Balances: bs,
	}, nil
}

type monthlyRecurringCost struct {
	name string
	date int
	currency.Code
	amount int
}

func newMonthlyRecurringCost(name string, date int, cc currency.Code, amount int) (*monthlyRecurringCost, error) {
	if date > maxMonthlyDate {
		return nil, fmt.Errorf("date cannot be more than %d", maxMonthlyDate)
	}
	return &monthlyRecurringCost{
		name:   name,
		date:   date,
		Code:   cc,
		amount: amount,
	}, nil
}

func (mrc monthlyRecurringCost) generateAccountBalances(times []time.Time) (AccountBalances, error) {
	// for each time, the balance is equal to the number of the specific dates that have passed
	return AccountBalances{}, errors.New("not implemented")
}
