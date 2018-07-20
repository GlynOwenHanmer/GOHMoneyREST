package main

import (
	"fmt"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	"github.com/pkg/errors"
)

type recurringCost interface {
	generateAccountBalances(times []time.Time) (AccountBalances, error)
}

const maxMonthlyDate = 28

type dailyRecurringCost struct {
	name string
	currency.Code
	Amount int
}

func (rcs dailyRecurringCost) generateAccountBalances(times []time.Time) (AccountBalances, error) {
	a, err := account.New(rcs.name, rcs.Code, time.Time{}) // time/date of account is not used currently
	if err != nil {
		return AccountBalances{}, errors.Wrap(err, "creating new account")
	}
	var bs balance.Balances
	for _, t := range times {
		var amount int
		if t.After(now) {
			amount = int(t.Sub(now)/(time.Hour*24)) * rcs.Amount
		}
		b, err := balance.New(t, balance.Amount(amount))
		if err != nil {
			return AccountBalances{}, errors.Wrap(err, "creating balance")
		}
		bs = append(bs, *b)
	}
	return AccountBalances{
		Account:  *a,
		Balances: bs,
	}, nil

	//for _, t := range times {
	//	// only occur cost if time is past now
	//
	//}
	//return nil, nil
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
	return AccountBalances{}, errors.New("not implemented")
}
