package date

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlag_SetErrors(t *testing.T) {
	for _, test := range []struct {
		name string
		vals []string
		time.Time
	}{
		{
			name: "zero-values",
		},
		{
			name: "whitespace",
			vals: []string{"\t\t\t\r"},
		},
		{
			name: "nonsense",
			vals: []string{"bloopy bleep", "!!!!!"},
		},
		{
			name: "invalid date format",
			vals: []string{"02/03", "-1000-01-87"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for key, val := range test.vals {
				k := key
				v := val
				t.Run(fmt.Sprintf("%d-%s", k, val), func(t *testing.T) {
					f := &flag{}
					err := f.Set(v)
					assert.Nil(t, f.Time)
					assert.Error(t, err)
				})
			}
		})
	}
}

func TestFlag_SetLongFormat(t *testing.T) {
	for _, test := range []struct {
		name string
		vals []string
		time.Time
	}{
		{
			name: "valid long format",
			vals: []string{"2018-03-02", "2018-3-2"},
			Time: time.Date(2018, 03, 02, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for key, val := range test.vals {
				k := key
				v := val
				t.Run(fmt.Sprintf("%d-%s", k, val), func(t *testing.T) {
					f := &flag{}
					err := f.Set(v)
					assert.NoError(t, err)
					if assert.NotNil(t, f) &&
						assert.NotNil(t, *f.Time) {
						assert.Equal(t, test.Time, *f.Time)
					}
				})
			}
		})
	}
}

func TestFlag_SetMonthDate(t *testing.T) {
	for _, test := range []struct {
		name string
		vals []string
		time.Time
	}{
		{
			name: "valid month date",
			vals: []string{"03-02", "3-2"},
			Time: time.Date(time.Now().Year(), 03, 02, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for key, val := range test.vals {
				k := key
				v := val
				t.Run(fmt.Sprintf("%d-%s", k, val), func(t *testing.T) {
					f := &flag{}
					err := f.Set(v)
					assert.NoError(t, err)
					if assert.NotNil(t, f) &&
						assert.NotNil(t, *f.Time) {
						assert.Equal(t, test.Time, *f.Time)
					}
				})
			}
		})
	}
}

func TestFlag_SetRelative(t *testing.T) {
	for _, test := range []struct {
		name         string
		vals         []string
		relativeDays int
	}{
		{
			name:         "yesterday",
			vals:         []string{"yesterday", "YESTERDAY", "y", "Y"},
			relativeDays: -1,
		},
		{
			name:         "relative date negative",
			vals:         []string{"-2"},
			relativeDays: -2,
		},
		{
			name:         "relative date positive",
			vals:         []string{"378", "+378"},
			relativeDays: 378,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for key, val := range test.vals {
				k := key
				v := val
				t.Run(fmt.Sprintf("%d-%s", k, val), func(t *testing.T) {
					f := &flag{}
					expected := time.Now().Add(time.Hour * 24 * time.Duration(test.relativeDays))
					err := f.Set(v)
					assert.NoError(t, err)
					assert.NotNil(t, f)
					diff := absDuration(f.Time.Sub(expected))
					acceptableThreshold := time.Millisecond * 10
					assert.Truef(t,
						diff < acceptableThreshold,
						"difference should be small but was %s. "+
							"If the difference is still quite small, "+
							"this could be because of a slow running test.",
						diff)
				})
			}
		})
	}
}

// absDuration will make any Duration positive.
// BUG: If the duration given is equal to math.MinInt64, which is a large
// negative number, the duration returned will be negative, still.
func absDuration(a time.Duration) time.Duration {
	if a >= 0 {
		return a
	}
	return a * -1
}
