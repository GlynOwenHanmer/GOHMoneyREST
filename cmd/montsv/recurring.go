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

type amountGenerator interface {
	generateAmount(time time.Time) int
	generateAccountBalances(times []time.Time) (AccountBalances, error)
}

type dailyRecurringCost struct {
	name string
	currency.Code
	Amount int
	from   time.Time
}

func (rcs dailyRecurringCost) generateAmount(at time.Time) int {
	var amount int
	if at.After(rcs.from) {
		amount = int(at.Sub(rcs.from)/(time.Hour*24)) * rcs.Amount
	}
	return amount
}

func (rcs dailyRecurringCost) generateBalance(at time.Time) (*balance.Balance, error) {
	b, err := balance.New(at, balance.Amount(rcs.generateAmount(at)))
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
	name        string
	from        time.Time
	dateOfMonth int
	amount      int
}

func newMonthlyRecurringCost(name string, dateOfMonth int, amount int) (*monthlyRecurringCost, error) {
	if dateOfMonth > maxMonthlyDate {
		return nil, fmt.Errorf("dateOfMonth cannot be more than %d", maxMonthlyDate)
	}
	return &monthlyRecurringCost{
		name:        name,
		dateOfMonth: dateOfMonth,
		amount:      amount,
	}, nil
}

func (mrc monthlyRecurringCost) generateAmount(at time.Time) int {
	offsetFromMonths := int(at.Month() - mrc.from.Month())
	var offsetFromDateOfMonth int
	//if at.Day() > mrc.from.Day() {
	//	offsetFromDateOfMonth = 1
	//}
	occurrences := offsetFromDateOfMonth + offsetFromMonths

	// for each at, the balance is equal to the number of the specific dates that have passed
	return occurrences * mrc.amount
}
