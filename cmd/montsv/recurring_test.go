package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonthlyRecurringCost(t *testing.T) {
	t.Run("should be zero at before from date", func(t *testing.T) {
		from := time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC)
		at := time.Date(500, 1, 1, 0, 0, 0, 0, time.UTC)

		mrc := monthlyRecurringCost{
			dateOfMonth: 1,
			from:        from,
		}

		expected := 0

		actual := mrc.generateAmount(at)
		assert.Equal(t, expected, actual)
	})

	t.Run("should be zero at from date", func(t *testing.T) {
		date := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)

		mrc := monthlyRecurringCost{
			dateOfMonth: 1,
			from:        date,
		}

		expected := 0

		actual := mrc.generateAmount(date)
		assert.Equal(t, expected, actual)
	})

	t.Run("should be zero at a day after from date when from.Day is equal to dateOfMonth", func(t *testing.T) {
		from := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		dateOfMonth := 1
		at := from.Add(time.Hour * 24)

		mrc := monthlyRecurringCost{
			dateOfMonth: dateOfMonth,
			from:        from,
			amount:      100,
		}

		expected := 0

		actual := mrc.generateAmount(at)
		assert.Equal(t, expected, actual)
	})

	t.Run("should be value of 1 occurrences at a month after from date", func(t *testing.T) {
		from := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		at := time.Date(0, 2, 1, 0, 0, 0, 0, time.UTC)

		mrc := monthlyRecurringCost{
			dateOfMonth: 1,
			from:        from,
			amount:      100,
		}

		expected := 100

		actual := mrc.generateAmount(at)
		assert.Equal(t, expected, actual, "%s %s", from.Format("Jan"), at.Format("Jan"))
	})

	t.Run("should be value of 1 occurrences at a day and a month after from date", func(t *testing.T) {
		from := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		at := time.Date(0, 2, 1, 0, 0, 0, 0, time.UTC)

		mrc := monthlyRecurringCost{
			dateOfMonth: 1,
			from:        from,
			amount:      100,
		}

		expected := 100

		actual := mrc.generateAmount(at)
		assert.Equal(t, expected, actual)
	})

	t.Run("should be value of 1 occurrence when from.Day is before dateOfMonth, at a day before dateOfMonth in the month after from date", func(t *testing.T) {
		from := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		dateOfMonth := 3
		at := time.Date(0, 2, 1, 0, 0, 0, 0, time.UTC)

		mrc := monthlyRecurringCost{
			dateOfMonth: dateOfMonth,
			from:        from,
			amount:      100,
		}

		expected := 100

		actual := mrc.generateAmount(at)
		assert.Equal(t, expected, actual)
	})

	t.Run("should be value of 2 occurrence when from is before dateOfMonth, at is after dateOfMonth and at is in next month", func(t *testing.T) {
		from := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		dateOfMonth := 2
		at := time.Date(0, 2, 3, 0, 0, 0, 0, time.UTC)

		mrc := monthlyRecurringCost{
			dateOfMonth: dateOfMonth,
			from:        from,
			amount:      100,
		}

		expected := 200

		actual := mrc.generateAmount(at)
		assert.Equal(t, expected, actual)
	})

}
