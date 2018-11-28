package main

import (
	"fmt"
	"time"
)

const maxMonthlyDate = 28

// amountGenerator should return the balance for a given time
// at current time of writing, all of the amountGenerators should return values
// only that represent the whole day. I.e. for a cost that occurs monthly on
// the 15th of a month, the value should be the same for any time given on the
// dat of the 15th of any month
type amountGenerator interface {
	generateAmount(time time.Time) int
}

type dailyRecurringAmount struct {
	Amount int
}

// generateAmount will return the amount of the daily recurring cost for any
// time that is passed.
func (dra dailyRecurringAmount) generateAmount(at time.Time) int {
	return dra.Amount
}

type monthlyRecurringCost struct {
	dateOfMonth int
	amount      int
}

func newMonthlyRecurringCost(dateOfMonth int, amount int) (*monthlyRecurringCost, error) {
	if dateOfMonth > maxMonthlyDate {
		return nil, fmt.Errorf("dateOfMonth cannot be more than %d", maxMonthlyDate)
	}
	return &monthlyRecurringCost{
		dateOfMonth: dateOfMonth,
		amount:      amount,
	}, nil
}

func (mrc monthlyRecurringCost) generateAmount(at time.Time) int {
	if mrc.dateOfMonth != at.Day() {
		return 0
	}
	return mrc.amount
}
