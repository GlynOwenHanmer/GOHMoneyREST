package main

import (
	"fmt"
	"time"
)

const maxMonthlyDate = 28

type amountGenerator interface {
	generateAmount(time time.Time) int
}

type dailyRecurringAmount struct {
	Amount int
	from   time.Time
}

func (rcs dailyRecurringAmount) generateAmount(at time.Time) int {
	var amount int
	if at.After(rcs.from) {
		amount = int(at.Sub(rcs.from)/(time.Hour*24)) * rcs.Amount
	}
	return amount
}

type monthlyRecurringCost struct {
	from        time.Time
	dateOfMonth int
	amount      int
}

func newMonthlyRecurringCost(name string, dateOfMonth int, amount int) (*monthlyRecurringCost, error) {
	if dateOfMonth > maxMonthlyDate {
		return nil, fmt.Errorf("dateOfMonth cannot be more than %d", maxMonthlyDate)
	}
	return &monthlyRecurringCost{
		dateOfMonth: dateOfMonth,
		amount:      amount,
	}, nil
}

func (mrc monthlyRecurringCost) generateAmount(at time.Time) int {
	offsetFromMonths := int(at.Month() - mrc.from.Month())
	var offsetFromDateOfMonth int
	if at.Day() > mrc.dateOfMonth {
		offsetFromDateOfMonth = 1
	}
	occurrences := offsetFromDateOfMonth + offsetFromMonths
	if occurrences < 0 {
		occurrences = 0
	}

	// for each at, the balance is equal to the number of the specific dates that have passed
	return occurrences * mrc.amount
}
